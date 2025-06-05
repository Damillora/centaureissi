use axum::{
    Json,
    extract::multipart::MultipartError,
    http::StatusCode,
    response::{IntoResponse, Response},
};
use serde::Serialize;
use tantivy::{TantivyError, query::QueryParserError};

pub enum CentaureissiError {
    Argon2Error(argon2::password_hash::Error),
    IncorrectPasswordError(),
    UserExistsError(),
    UserNotFoundError(),
    InternalServerError(String),
    AuthError(String, StatusCode),
    UnimplementedError(),
    InvalidDataError(MultipartError),
    InvalidEmailContentsError(String),
    RelationalDatabaseError(),
    BlobDatabaseError(),
    RegistrationDisabled(),
    MessageNotFound(),
    MessageError(),
}

impl IntoResponse for CentaureissiError {
    fn into_response(self) -> Response {
        // How we want errors responses to be serialized
        #[derive(Serialize)]
        struct ErrorResponse {
            message: String,
        }

        let (status, message) = match self {
            Self::Argon2Error(err) => (StatusCode::BAD_REQUEST, err.to_string()),
            Self::IncorrectPasswordError() => (
                StatusCode::UNAUTHORIZED,
                "Incorrect username or password".to_string(),
            ),
            Self::UserExistsError() => (StatusCode::BAD_REQUEST, "User already exists".to_string()),
            Self::UserNotFoundError() => (StatusCode::BAD_REQUEST, "User not found".to_string()),
            Self::InternalServerError(err) => (
                StatusCode::INTERNAL_SERVER_ERROR,
                format!("Something happened: {}", err),
            ),
            Self::UnimplementedError() => (
                StatusCode::INTERNAL_SERVER_ERROR,
                "API unimplemented!".to_string(),
            ),
            Self::AuthError(message, status) => (status, message),
            Self::InvalidDataError(err) => (
                StatusCode::BAD_REQUEST,
                format!("Invalid email data: {}", err.to_string()),
            ),
            Self::InvalidEmailContentsError(err) => (
                StatusCode::BAD_REQUEST,
                format!("Invalid email contents: {}", err),
            ),
            Self::RelationalDatabaseError() => (
                StatusCode::INTERNAL_SERVER_ERROR,
                "Relational database error".to_string(),
            ),
            Self::BlobDatabaseError() => (
                StatusCode::INTERNAL_SERVER_ERROR,
                "Blob database error".to_string(),
            ),
            Self::RegistrationDisabled() => (
                StatusCode::BAD_REQUEST,
                "Registration is disabled".to_string(),
            ),
            Self::MessageNotFound() => (StatusCode::NOT_FOUND, "Message not found".to_string()),
            Self::MessageError() => (StatusCode::BAD_REQUEST, "Message is unreadable".to_string()),
        };

        (status, Json(ErrorResponse { message })).into_response()
    }
}

impl From<argon2::password_hash::Error> for CentaureissiError {
    fn from(error: argon2::password_hash::Error) -> Self {
        Self::Argon2Error(error)
    }
}

impl<T: Into<persy::PersyError>> From<persy::PE<T>> for CentaureissiError {
    fn from(err: persy::PE<T>) -> Self {
        Self::InternalServerError(format!(
            "Blob storage error: {}",
            err.persy_error().to_string()
        ))
    }
}

impl From<TantivyError> for CentaureissiError {
    fn from(err: TantivyError) -> Self {
        Self::InternalServerError(format!("Search error: {}", err.to_string()))
    }
}

impl From<QueryParserError> for CentaureissiError {
    fn from(err: QueryParserError) -> Self {
        Self::InternalServerError(format!("Search error: {}", err.to_string()))
    }
}
