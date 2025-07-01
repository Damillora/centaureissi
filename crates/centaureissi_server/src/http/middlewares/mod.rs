use std::sync::Arc;

use axum::{
    body::Body,
    extract::{Query, Request, State},
    http::{self, StatusCode},
    middleware::Next,
    response::Response,
};
use diesel::{QueryDsl, RunQueryDsl, SelectableHelper};
use jsonwebtoken::{DecodingKey, TokenData, Validation, decode};
use tracing::info;

use crate::{db::models::User, http::models::requests::DownloadTokenRequest};

use super::{context::CentaureissiContext, errors::CentaureissiError, models::auth::Claims};

pub async fn authorization_middleware(
    State(context): State<Arc<CentaureissiContext>>,
    mut req: Request,
    next: Next,
) -> Result<Response<Body>, CentaureissiError> {
    let auth_header = req.headers().get(http::header::AUTHORIZATION);
    let auth_header = match auth_header {
        Some(header) => header.to_str().map_err(|_| {
            CentaureissiError::AuthError(
                "Empty header is not allowed".to_string(),
                StatusCode::FORBIDDEN,
            )
        })?,
        None => {
            return Err(CentaureissiError::AuthError(
                "Please add the JWT token to the header".to_string(),
                StatusCode::FORBIDDEN,
            ));
        }
    };
    let mut header = auth_header.split_whitespace();
    let (_, token) = (header.next(), header.next());

    let secret = context.config.auth_secret.clone();
    let mut validator = Validation::default();
    validator.set_audience(&["centaureissi"]);
    let token_data: TokenData<Claims> = decode(
        &token.unwrap().to_string(),
        &DecodingKey::from_secret(secret.as_ref()),
        &validator,
    )
    .map_err(|_| {
        CentaureissiError::AuthError(
            "Unable to parse token".to_string(),
            StatusCode::UNAUTHORIZED,
        )
    })?;
    {
        use crate::db::schema::users::dsl::*;

        let conn = &mut context.db.get().unwrap();

        let user = users
            .find(token_data.claims.sub)
            .select(User::as_select())
            .first(conn)
            .map_err(|_| CentaureissiError::IncorrectPasswordError())?;

        req.extensions_mut().insert(user);
    }
    Ok(next.run(req).await)
}

pub async fn authorization_sign_middleware(
    State(context): State<Arc<CentaureissiContext>>,
    Query(input): Query<DownloadTokenRequest>,
    mut req: Request,
    next: Next,
) -> Result<Response<Body>, CentaureissiError> {
    let token = input.token;
    let secret = context.config.auth_secret.clone();
    let mut validator = Validation::default();
    validator.set_audience(&[req.uri().path()]);
    let token_data: TokenData<Claims> = decode(
        &token.to_string(),
        &DecodingKey::from_secret(secret.as_ref()),
        &validator,
    )
    .map_err(|_| {
        CentaureissiError::AuthError(
            "Unable to parse token".to_string(),
            StatusCode::UNAUTHORIZED,
        )
    })?;
    {
        use crate::db::schema::users::dsl::*;

        let conn = &mut context.db.get().unwrap();

        let user = users
            .find(token_data.claims.sub)
            .select(User::as_select())
            .first(conn)
            .map_err(|_| CentaureissiError::IncorrectPasswordError())?;

        req.extensions_mut().insert(user);
    }
    Ok(next.run(req).await)
}
