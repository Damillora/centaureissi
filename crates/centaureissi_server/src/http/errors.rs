use axum::{http::StatusCode, response::{IntoResponse, Response}, Json};
use serde::Serialize;
pub enum CentaureissiError {
    Argon2Error(argon2::password_hash::Error),
    IncorrectPasswordError(),
    UserExistsError(),
    UserNotFoundError(),
    InternalServerError(),
    AuthError(String, StatusCode),
    Unimplemented(),
}

impl IntoResponse for CentaureissiError {
    fn into_response(self) -> Response {
        // How we want errors responses to be serialized
        #[derive(Serialize)]
        struct ErrorResponse {
            message: String,
        }

        let (status, message) = match self {
            Self::Argon2Error(err) => {
                (StatusCode::BAD_REQUEST, err.to_string())
            }
            Self::IncorrectPasswordError() => {
                (StatusCode::UNAUTHORIZED, "Incorrect username or password".to_string())
            },
            Self::UserExistsError() => {
                (StatusCode::BAD_REQUEST,  "User already exists".to_string())
            }
            Self::UserNotFoundError() => {
                (StatusCode::BAD_REQUEST,  "User not found".to_string())
            }
            Self::InternalServerError() => {
                (StatusCode::INTERNAL_SERVER_ERROR, "Something happened!".to_string())
            }
            Self::Unimplemented() => {
                (StatusCode::INTERNAL_SERVER_ERROR, "API unimplemented!".to_string())
            }
            Self::AuthError(message,  status) => {
                (status, message)
            }
        };

        (status, Json(ErrorResponse { message })).into_response()
    }
}

impl From<argon2::password_hash::Error> for CentaureissiError {
    fn from(error: argon2::password_hash::Error) -> Self {
        Self::Argon2Error(error)
    }
}