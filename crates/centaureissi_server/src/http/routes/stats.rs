use std::sync::Arc;

use axum::{Json, Router, extract::State, middleware, routing::get};
use diesel::{QueryDsl, RunQueryDsl};
use tokio::fs;

use crate::{
    blobs::BLOB_TABLE,
    http::{
        context::CentaureissiContext, errors::CentaureissiError, middlewares,
        models::responses::StatsResponse,
    },
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new();

    let protected_router =
        Router::new()
            .route("/", get(stats))
            .layer(middleware::from_fn_with_state(
                state,
                middlewares::authorization_middleware,
            ));

    return unprotected_router.merge(protected_router);
}

async fn stats(
    State(context): State<Arc<CentaureissiContext>>,
) -> Result<Json<StatsResponse>, CentaureissiError> {
    let database_url = context.config.get_database_url();
    let blob_db_file = context.config.get_blob_db_path();

    let db_size = fs::metadata(database_url).await.unwrap().len();
    let blob_db_size = fs::metadata(blob_db_file).await.unwrap().len();

    use crate::db::schema::messages::dsl::*;

    let conn = &mut context.db.get().unwrap();

    let message_count: i64 = messages
        .count()
        .get_result(conn)
        .map_err(|_| CentaureissiError::RelationalDatabaseError())?;
    let searcher = context.search_reader.searcher();
    let search_doc_count = searcher.num_docs();
    let mut blob_db_count: u64 = 0;
    for _ in context.blob_db.scan(BLOB_TABLE)? {
        blob_db_count += 1;
    }

    let response = StatsResponse {
        version: env!("CARGO_PKG_VERSION").to_string(),
        db_size: db_size,
        message_count: message_count as u64,
        blob_db_size: blob_db_size,
        blob_count: blob_db_count as u64,
        search_doc_count: search_doc_count,
    };
    Ok(Json(response))
}
