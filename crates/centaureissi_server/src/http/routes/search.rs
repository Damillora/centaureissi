use std::sync::Arc;

use axum::{
    Extension, Json, Router,
    extract::{Query, State},
    middleware,
    routing::get,
};
use chrono::NaiveDateTime;
use mail_parser::DateTime;
use tantivy::{
    Document, TantivyDocument,
    collector::{Count, TopDocs},
    query::{self, QueryParser},
    schema::Value,
};

use crate::{
    db::models::User,
    http::{
        context::CentaureissiContext,
        errors::CentaureissiError,
        middlewares,
        models::{
            requests::SearchRequest,
            responses::{SearchResponse, SearchResponseItem},
        },
    },
    search,
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new();

    let protected_router =
        Router::new()
            .route("/", get(search_message))
            .layer(middleware::from_fn_with_state(
                state,
                middlewares::authorization_middleware,
            ));

    return unprotected_router.merge(protected_router);
}

async fn search_message(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Query(input): Query<SearchRequest>,
) -> Result<Json<SearchResponse>, CentaureissiError> {
    let searcher = context.search_reader.searcher();

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

    let query_parser = QueryParser::for_index(
        &context.search_index,
        vec![from, to, cc, bcc, subject, content],
    );
    let query = query_parser.parse_query(&input.q)?;

    let mut page = input.page.unwrap_or(1);
    if page < 1 {
        page = 1;
    }
    let mut per_page = input.per_page.unwrap_or(10);
    if per_page < 1 {
        per_page = 10;
    }

    let offset = per_page * (page - 1);

    let total_hits = searcher.search(&query, &Count)?;
    let mut total_pages = total_hits / per_page;
    if total_hits % per_page != 0 {
        total_pages = total_pages + 1;
    }

    let result = searcher
        .search(&query, &TopDocs::with_limit(10).and_offset(offset))?
        .into_iter()
        .map(|(_score, doc_address)| {
            let doc = searcher.doc::<TantivyDocument>(doc_address).unwrap();

            // if context.config.verbose {

            //     let json: Vec<String> = doc.iter_fields_and_values()
            //         .map(|(field, item)| format!("{} {}", field.field_id(), item.is_null()))
            //         .collect();

            //     println!("Fields:\n{}", json.join("\n"))
            // }

            SearchResponseItem {
                id: doc.get_first(id).unwrap().as_i64().unwrap(),
                hash: doc.get_first(hash).unwrap().as_str().unwrap().to_string(),
                user_id: doc.get_first(user_id).unwrap().as_i64().unwrap(),

                from: doc.get_first(from).unwrap().as_str().unwrap().to_string(),
                to: doc.get_first(to).unwrap().as_str().unwrap().to_string(),
                cc: doc.get_first(cc).unwrap().as_str().unwrap().to_string(),
                bcc: doc.get_first(bcc).unwrap().as_str().unwrap().to_string(),
                subject: doc
                    .get_first(subject)
                    .unwrap()
                    .as_str()
                    .unwrap()
                    .to_string(),
                date: DateTime::from_timestamp(
                    doc
                    .get_first(date)
                    .unwrap()
                    .as_datetime()
                    .unwrap()
                    .into_timestamp_secs()
                ).to_rfc3339(),
            }
        });

    let response = SearchResponse {
        items: result.collect(),
        page: page,
        total_items: total_hits,
        total_pages: total_pages,
    };
    Ok(Json(response))
}
