use std::sync::Arc;

use argon2::{Argon2, PasswordHash, PasswordVerifier};
use axum::{Extension, Json, Router, extract::State, http::StatusCode, middleware, routing::post};
use chrono::{Duration, Utc};
use diesel::{
    ExpressionMethods, RunQueryDsl, SelectableHelper,
    query_dsl::methods::{FilterDsl, SelectDsl},
};
use jsonwebtoken::{EncodingKey, Header, encode};

use crate::{
    db::models::User,
    http::{
        context::CentaureissiContext,
        errors::CentaureissiError,
        middlewares,
        models::{
            auth::Claims,
            requests::{LoginRequest, SignPathRequest},
            responses::LoginResponse,
        },
    },
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new().route("/login", post(login));

    let protected_router = Router::new()
        .route("/token", post(get_token))
        .route("/token/sign", post(get_token_sign))
        .layer(middleware::from_fn_with_state(
            state,
            middlewares::authorization_middleware,
        ));

    return unprotected_router.merge(protected_router);
}

async fn login(
    State(context): State<Arc<CentaureissiContext>>,
    Json(input): Json<LoginRequest>,
) -> Result<Json<LoginResponse>, CentaureissiError> {
    use crate::db::schema::users::dsl::*;

    let conn = &mut context.db.get().unwrap();

    let user = users
        .filter(username.eq(input.username))
        .select(User::as_select())
        .first(conn)
        .map_err(|_| CentaureissiError::IncorrectPasswordError())?;

    let hashed_password = PasswordHash::new(&user.password)?;
    let argon2 = Argon2::default();

    argon2
        .verify_password(input.password.as_bytes(), &hashed_password)
        .map_err(|_| CentaureissiError::IncorrectPasswordError())?;

    let secret = context.config.auth_secret.clone();
    let now = Utc::now();
    let expire: chrono::TimeDelta = Duration::hours(24);
    let exp: usize = (now + expire).timestamp() as usize;
    let iat: usize = now.timestamp() as usize;
    let claims = Claims {
        aud: "centaureissi".to_string(),
        iss: "centaureissi-api".to_string(),
        exp: exp,
        iat: iat,
        name: user.username,
        sub: user.id,
    };

    let token = encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_ref()),
    )
    .map_err(|_| {
        CentaureissiError::AuthError(
            String::from("Failure to decode token"),
            StatusCode::FORBIDDEN,
        )
    })?;

    Ok(Json(LoginResponse { token: token }))
}

async fn get_token(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
) -> Result<Json<LoginResponse>, CentaureissiError> {
    let secret = context.config.auth_secret.clone();
    let now = Utc::now();
    let expire: chrono::TimeDelta = Duration::hours(24);
    let exp: usize = (now + expire).timestamp() as usize;
    let iat: usize = now.timestamp() as usize;
    let claims = Claims {
        aud: "centaureissi".to_string(),
        iss: "centaureissi-api".to_string(),
        exp: exp,
        iat: iat,
        name: user.username,
        sub: user.id,
    };

    let token = encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_ref()),
    )
    .map_err(|_| {
        CentaureissiError::AuthError(
            String::from("Failure to decode token"),
            StatusCode::FORBIDDEN,
        )
    })?;

    Ok(Json(LoginResponse { token: token }))
}

async fn get_token_sign(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Json(input): Json<SignPathRequest>,
) -> Result<Json<LoginResponse>, CentaureissiError> {
    let secret = context.config.auth_secret.clone();
    let now = Utc::now();
    let expire: chrono::TimeDelta = Duration::hours(24);
    let exp: usize = (now + expire).timestamp() as usize;
    let iat: usize = now.timestamp() as usize;
    let claims = Claims {
        aud: input.path.to_string(),
        iss: "centaureissi-api".to_string(),
        exp: exp,
        iat: iat,
        name: user.username,
        sub: user.id,
    };

    let token = encode(
        &Header::default(),
        &claims,
        &EncodingKey::from_secret(secret.as_ref()),
    )
    .map_err(|_| {
        CentaureissiError::AuthError(
            String::from("Failure to decode token"),
            StatusCode::FORBIDDEN,
        )
    })?;

    Ok(Json(LoginResponse { token: token }))
}
