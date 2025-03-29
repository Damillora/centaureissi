use axum::{
    Router,
    http::{StatusCode, Uri, header},
    response::IntoResponse,
    routing::get,
};

use crate::http::context::CentaureissiContext;
use centaureissi_web::WebEmbed;

pub fn router() -> Router<CentaureissiContext> {
    return Router::new()
        .route("/", get(index_handler))
        .route("/index.html", get(index_handler))
        .route("/_app/{*file}", get(static_handler))
        .fallback_service(get(index_handler));
}

// We use static route matchers ("/" and "/index.html") to serve our home
// page.
async fn index_handler() -> impl IntoResponse {
    static_handler("/app.html".parse::<Uri>().unwrap()).await
}

// We use a wildcard matcher ("/dist/*file") to match against everything
// within our defined assets directory. This is the directory on our Asset
// struct below, where folder = "examples/public/".
async fn static_handler(uri: Uri) -> impl IntoResponse {
    let mut path = uri.path().trim_start_matches('/').to_string();

    match WebEmbed::get(path.as_str()) {
        Some(content) => {
            let mime = mime_guess::from_path(path).first_or_octet_stream();
            ([(header::CONTENT_TYPE, mime.as_ref())], content.data).into_response()
        }
        None => (StatusCode::NOT_FOUND, "404 Not Found").into_response(),
    }
}
