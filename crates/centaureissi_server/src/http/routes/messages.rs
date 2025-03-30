use std::sync::Arc;

use axum::{
    Extension, Router,
    extract::{Multipart, State},
    http::StatusCode,
    middleware,
    response::IntoResponse,
    routing::post,
};
use blake2::{Blake2b512, Digest};
use diesel::{RunQueryDsl, SelectableHelper};
use mail_parser::MessageParser;
use persy::PersyId;
use tantivy::doc;

use crate::{
    blobs::{BLOB_INDEX, BLOB_TABLE},
    db::{
        models::{Messages, User},
        requests::NewMessage,
    },
    http::{context::CentaureissiContext, errors::CentaureissiError, middlewares},
    search,
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new();

    let protected_router =
        Router::new()
            .route("/add", post(index_message))
            .layer(middleware::from_fn_with_state(
                state,
                middlewares::authorization_middleware,
            ));

    return unprotected_router.merge(protected_router);
}

async fn index_message(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    mut multipart: Multipart,
) -> Result<impl IntoResponse, CentaureissiError> {
    let conn = &mut context.db.get().unwrap();
    use crate::db::schema::messages;

    while let Some(field) = multipart.next_field().await.unwrap() {
        let file_name = field.file_name().unwrap().to_string();
        let data = field
            .bytes()
            .await
            .map_err(|err| CentaureissiError::InvalidDataError(err))?;

        if context.config.verbose {
            println!("`{file_name}`:  {} bytes", data.len());
        }

        // Get content hash
        let mut hasher = Blake2b512::new();
        hasher.update(&data);
        let hash_array = hasher.finalize();
        let content_hash = format!("{:x}", hash_array);
        if context.config.verbose {
            println!("`{file_name}: content hash {}", &content_hash,)
        }

        // Check existence of existing blobs
        let mut blob_item = context
            .blob_db
            .get::<String, PersyId>(BLOB_INDEX, &content_hash.to_string())?;
        if let Some(_) = blob_item.next() {
            // Already in blob database, no need to add
            continue;
        }

        // ZSTD compression
        let mut compressed_data = Vec::<u8>::new();
        zstd::stream::copy_encode(&*data, &mut compressed_data, 0)
            .map_err(|_| CentaureissiError::BlobDatabaseError())?;

        // Insert into database
        let new_message = NewMessage {
            user_id: user.id,
            content_hash: content_hash.clone(),
        };

        let message = diesel::insert_into(messages::table)
            .values(&new_message)
            .returning(Messages::as_returning())
            .get_result(conn)
            .map_err(|_| CentaureissiError::RelationalDatabaseError())?;

        // Insert email to blob transaction
        let mut blob_write_txn = context.blob_db.begin()?;
        let id = blob_write_txn.insert(BLOB_TABLE, &compressed_data)?;
        blob_write_txn.put::<String, PersyId>(BLOB_INDEX, content_hash.clone(), id)?;
        let prepared = blob_write_txn.prepare()?;
        prepared.commit()?;

        // Parse email for search indexing
        let data_vec = data.to_vec();
        let parsed_msg = MessageParser::default().parse(&data_vec);
        if let Some(msg) = parsed_msg {
            let from_data: Vec<String> = match msg.from() {
                Some(from) => from
                    .iter()
                    .map(|item| {
                        if let Some(name) = item.name() {
                            format!("{} <{}>", name, item.address().unwrap())
                        } else {
                            format!("{}", item.address().unwrap())
                        }
                    })
                    .collect(),
                None => Vec::new(),
            };
            let to_data: Vec<String> = match msg.to() {
                Some(to) => to
                    .iter()
                    .map(|item| {
                        if let Some(name) = item.name() {
                            format!("{} <{}>", name, item.address().unwrap_or(""))
                        } else {
                            format!("{}", item.address().unwrap_or(""))
                        }
                    })
                    .collect(),
                None => Vec::new(),
            };
            let cc_data: Vec<String> = match msg.cc() {
                Some(cc) => cc
                    .iter()
                    .map(|item| {
                        if let Some(name) = item.name() {
                            format!("{} <{}>", name, item.address().unwrap())
                        } else {
                            format!("{}", item.address().unwrap())
                        }
                    })
                    .collect(),
                None => Vec::new(),
            };
            let bcc_data: Vec<String> = match msg.bcc() {
                Some(bcc) => bcc
                    .iter()
                    .map(|item| {
                        if let Some(name) = item.name() {
                            format!("{} <{}>", name, item.address().unwrap())
                        } else {
                            format!("{}", item.address().unwrap())
                        }
                    })
                    .collect(),
                None => Vec::new(),
            };
            let subject_data = msg.subject().unwrap_or("");
            let date_data = msg
                .date()
                .unwrap_or(&mail_parser::DateTime::from_timestamp(0))
                .to_timestamp();
            let mail_contents_data: Vec<String> = msg
                .text_bodies()
                .map(|item| item.text_contents().unwrap_or("").to_string())
                .collect();

            let schema = search::get_schema();

            // Schema Fields
            let id = schema.get_field("id").unwrap();
            let hash = schema.get_field("hash").unwrap();
            let user_id = schema.get_field("user_id").unwrap();
            let from = schema.get_field("from").unwrap();
            let to = schema.get_field("to").unwrap();
            let cc = schema.get_field("cc").unwrap();
            let bcc = schema.get_field("bcc").unwrap();
            let subject = schema.get_field("subject").unwrap();
            let date = schema.get_field("date").unwrap();
            let content = schema.get_field("content").unwrap();

            let mut search_adder = context.search_writer.write().unwrap();
            let search_doc = doc!(
                id => i64::from(message.id),
                hash => content_hash,
                user_id => i64::from(user.id),
                from => from_data.join(", "),
                to => to_data.join(", "),
                cc => cc_data.join(", "),
                bcc => bcc_data.join(", "),
                subject => subject_data,
                date => tantivy::DateTime::from_timestamp_secs(date_data),
                content => mail_contents_data.join("\n\n"),
            );
            search_adder.add_document(search_doc)?;
            search_adder.commit()?;
        } else {
            return Err(CentaureissiError::InvalidEmailContentsError(String::from(
                "no message",
            )));
        }
    }

    Ok((StatusCode::OK, ""))
}
