use std::{net::SocketAddr, sync::Arc};

use axum::Router;
use context::CentaureissiContext;
use diesel::{
    SqliteConnection,
    r2d2::{ConnectionManager, Pool},
};

use crate::config::CentaureissiConfig;

pub mod context;
pub mod errors;
pub mod middlewares;
pub mod models;
pub mod routes;

pub async fn serve(config: CentaureissiConfig, db: Pool<ConnectionManager<SqliteConnection>>) {
    let context = CentaureissiContext {
        config: Arc::new(config),
        db: db,
    };

    let app = Router::new()
        .nest("/api/user", routes::users::router(context.clone()))
        .nest("/api/auth", routes::auth::router(context.clone()))
        .merge(routes::web::router())
        .with_state(context);

    // Listen
    let addr = SocketAddr::from(([127, 0, 0, 1], 8080));
    tracing::debug!("listening on {addr}");
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();

    // Go
    axum::serve(listener, app.into_make_service())
        .await
        .unwrap();
}
