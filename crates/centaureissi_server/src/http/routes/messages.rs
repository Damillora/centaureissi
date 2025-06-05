use std::sync::Arc;

use axum::{
    Extension, Json, Router,
    extract::{Multipart, Path, State},
    http::StatusCode,
    middleware,
    response::IntoResponse,
    routing::{get, post},
};
use blake2::{Blake2b512, Digest};
use chrono::DateTime;
use diesel::{RunQueryDsl, SelectableHelper};
use mail_parser::MessageParser;
use persy::PersyId;

use crate::{
    blobs::{BLOB_INDEX, BLOB_TABLE},
    db::{
        models::{Messages, User},
        requests::NewMessage,
    },
    http::{
        context::CentaureissiContext,
        errors::CentaureissiError,
        middlewares,
        models::responses::{MessageContentResponse, MessageResponse},
    },
    search, utils,
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new();

    let protected_router = Router::new()
        .route("/add", post(index_message))
        .route("/{message_hash}/info", get(read_message_info))
        .route("/{message_hash}/raw", get(read_message_raw))
        .route("/{message_hash}/text", get(read_message_text))
        .route("/{message_hash}/html", get(read_message_html))
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

        diesel::insert_into(messages::table)
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
            let search_doc =
                search::message::create_search_document_from_message(user.id, content_hash, msg);

            let mut search_adder = context.search_writer.write().unwrap();
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

async fn read_message_info(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Path(message_hash): Path<String>,
) -> Result<impl IntoResponse, CentaureissiError> {
    println!("{}", message_hash);

    let blob_id = context
        .blob_db
        .one::<String, PersyId>(BLOB_INDEX, &message_hash)
        .unwrap();

    if let Some(blob_id) = blob_id {
        let blob_contents = context.blob_db.read(BLOB_TABLE, &blob_id).unwrap();

        if let Some(contents) = blob_contents {
            let mut uncompressed = Vec::<u8>::new();

            let is_compressed = !zstd::stream::copy_decode(&*contents, &mut uncompressed).is_err();

            if !is_compressed {
                uncompressed = contents;
            }

            let parsed_msg = MessageParser::default().parse(&uncompressed);

            if let Some(msg) = parsed_msg {
                let message_doc = utils::message::create_message_model_from_message(user.id, msg);

                let date = DateTime::from_timestamp(message_doc.timestamp_secs, 0).unwrap();

                return Ok(Json(MessageResponse {
                    hash: message_hash,

                    date: date.to_rfc3339().to_string(),
                    from: message_doc.from,
                    to: message_doc.to,
                    cc: message_doc.cc,
                    bcc: message_doc.bcc,

                    is_html_mail: message_doc.is_html_mail,
                    is_text_mail: message_doc.is_text_mail,
                    subject: message_doc.subject,
                }));
            } else {
                return Err(CentaureissiError::MessageError());
            }
        } else {
            return Err(CentaureissiError::MessageNotFound());
        }
    } else {
        return Err(CentaureissiError::MessageNotFound());
    }
}

async fn read_message_raw(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Path(message_hash): Path<String>,
) -> Result<impl IntoResponse, CentaureissiError> {
    println!("{}", message_hash);

    let blob_id = context
        .blob_db
        .one::<String, PersyId>(BLOB_INDEX, &message_hash)
        .unwrap();

    if let Some(blob_id) = blob_id {
        let blob_contents = context.blob_db.read(BLOB_TABLE, &blob_id).unwrap();

        if let Some(contents) = blob_contents {
            let mut uncompressed = Vec::<u8>::new();

            let is_compressed = !zstd::stream::copy_decode(&*contents, &mut uncompressed).is_err();

            if !is_compressed {
                uncompressed = contents;
            }

            return Ok(Json(MessageContentResponse {
                content: String::from_utf8(uncompressed).unwrap(),
            }));
        } else {
            return Err(CentaureissiError::MessageNotFound());
        }
    } else {
        return Err(CentaureissiError::MessageNotFound());
    }
}

async fn read_message_html(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Path(message_hash): Path<String>,
) -> Result<impl IntoResponse, CentaureissiError> {
    println!("{}", message_hash);

    let blob_id = context
        .blob_db
        .one::<String, PersyId>(BLOB_INDEX, &message_hash)
        .unwrap();

    if let Some(blob_id) = blob_id {
        let blob_contents = context.blob_db.read(BLOB_TABLE, &blob_id).unwrap();

        if let Some(contents) = blob_contents {
            let mut uncompressed = Vec::<u8>::new();

            let is_compressed = !zstd::stream::copy_decode(&*contents, &mut uncompressed).is_err();

            if !is_compressed {
                uncompressed = contents;
            }

            let parsed_msg = MessageParser::default().parse(&uncompressed);

            if let Some(msg) = parsed_msg {
                let message_doc =
                    utils::message::create_message_content_model_from_message(true, msg);

                return Ok(Json(MessageContentResponse {
                    content: message_doc.content,
                }));
            } else {
                return Err(CentaureissiError::MessageError());
            }
        } else {
            return Err(CentaureissiError::MessageNotFound());
        }
    } else {
        return Err(CentaureissiError::MessageNotFound());
    }
}

async fn read_message_text(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Path(message_hash): Path<String>,
) -> Result<impl IntoResponse, CentaureissiError> {
    println!("{}", message_hash);

    let blob_id = context
        .blob_db
        .one::<String, PersyId>(BLOB_INDEX, &message_hash)
        .unwrap();

    if let Some(blob_id) = blob_id {
        let blob_contents = context.blob_db.read(BLOB_TABLE, &blob_id).unwrap();

        if let Some(contents) = blob_contents {
            let mut uncompressed = Vec::<u8>::new();

            let is_compressed = !zstd::stream::copy_decode(&*contents, &mut uncompressed).is_err();

            if !is_compressed {
                uncompressed = contents;
            }

            let parsed_msg = MessageParser::default().parse(&uncompressed);

            if let Some(msg) = parsed_msg {
                let message_doc =
                    utils::message::create_message_content_model_from_message(false, msg);

                return Ok(Json(MessageContentResponse {
                    content: message_doc.content,
                }));
            } else {
                return Err(CentaureissiError::MessageError());
            }
        } else {
            return Err(CentaureissiError::MessageNotFound());
        }
    } else {
        return Err(CentaureissiError::MessageNotFound());
    }
}
