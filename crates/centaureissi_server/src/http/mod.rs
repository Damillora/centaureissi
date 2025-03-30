use std::{
    net::SocketAddr,
    sync::{Arc, RwLock},
};

use axum::{Router, extract::DefaultBodyLimit};
use context::CentaureissiContext;
use diesel::{
    SqliteConnection,
    r2d2::{ConnectionManager, Pool},
};
use tantivy::{Index, IndexReader, IndexWriter};

use crate::config::CentaureissiConfig;

pub mod context;
pub mod errors;
pub mod middlewares;
pub mod models;
pub mod routes;

pub async fn serve(
    config: CentaureissiConfig,
    db: Pool<ConnectionManager<SqliteConnection>>,
    index: Index,
    index_writer: IndexWriter,
    index_reader: IndexReader,
    blob_db: persy::Persy,
) {
    let context = CentaureissiContext {
        config: config,
        db: db,
        search_index: index,
        search_writer: RwLock::new(index_writer),
        search_reader: index_reader,
        blob_db: blob_db,
    };
    let shared_context = Arc::new(context);

    let app = Router::new()
        .nest("/api/user", routes::users::router(shared_context.clone()))
        .nest("/api/auth", routes::auth::router(shared_context.clone()))
        .nest(
            "/api/messages",
            routes::messages::router(shared_context.clone()),
        )
        .nest(
            "/api/search",
            routes::search::router(shared_context.clone()),
        )
        .nest(
            "/api/stats",
            routes::stats::router(shared_context.clone()),
        )
        .merge(routes::web::router())
        // Replace the default of 2MB with 100MB
        .layer(DefaultBodyLimit::max(100_000_000))
        .with_state(shared_context);

    // Listen
    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    tracing::debug!("listening on {addr}");
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();

    // Go
    axum::serve(listener, app.into_make_service())
        .await
        .unwrap();
}
