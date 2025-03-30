use std::sync::Arc;

use argon2::{
    Argon2, PasswordHash, PasswordVerifier,
    password_hash::{PasswordHasher, SaltString, rand_core::OsRng},
};
use axum::{
    Extension, Json, Router,
    extract::State,
    http::StatusCode,
    middleware,
    response::IntoResponse,
    routing::{get, post, put},
};
use diesel::{RunQueryDsl, SelectableHelper};

use crate::{
    db::{
        self,
        models::User,
        requests::{UpdateUser, UpdateUserPassword},
        responses::UserProfile,
    },
    http::{
        context::CentaureissiContext,
        errors::CentaureissiError,
        middlewares,
        models::{
            requests::{NewUserRequest, UserUpdatePasswordRequest, UserUpdateRequest},
            responses::UserProfileResponse,
        },
    },
};

pub fn router(state: Arc<CentaureissiContext>) -> Router<Arc<CentaureissiContext>> {
    let unprotected_router = Router::new().route("/register", post(register_user));

    let protected_router = Router::new()
        .route("/profile", get(user_profile))
        .route("/update", put(user_update))
        .route("/update-password", put(user_update_password))
        .layer(middleware::from_fn_with_state(
            state,
            middlewares::authorization_middleware,
        ));

    return unprotected_router.merge(protected_router);
}

async fn register_user(
    State(context): State<Arc<CentaureissiContext>>,
    Json(input): Json<NewUserRequest>,
) -> Result<Json<UserProfileResponse>, CentaureissiError> {
    use crate::db::schema::users;
    if context.config.disable_registration {
        return Err(CentaureissiError::RegistrationDisabled());
    }
    let conn = &mut context.db.get().unwrap();

    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();
    let hashed_password = argon2
        .hash_password(input.password.as_bytes(), &salt)?
        .to_string();

    let new_user_request = db::requests::NewUser {
        username: input.username,
        password: hashed_password,
    };

    let user = diesel::insert_into(users::table)
        .values(&new_user_request)
        .returning(UserProfile::as_returning())
        .get_result(conn)
        .map_err(|err| match err {
            diesel::result::Error::DatabaseError(_, _) => CentaureissiError::UserExistsError(),
            _ => CentaureissiError::RelationalDatabaseError(),
        })?;

    Ok(Json(UserProfileResponse {
        username: user.username,
    }))
}

async fn user_profile(
    Extension(user): Extension<User>,
) -> Result<Json<UserProfileResponse>, CentaureissiError> {
    Ok(Json(UserProfileResponse {
        username: user.username,
    }))
}

async fn user_update(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Json(input): Json<UserUpdateRequest>,
) -> Result<impl IntoResponse, CentaureissiError> {
    let conn = &mut context.db.get().unwrap();

    let changeset = UpdateUser {
        username: input.username,
    };

    diesel::update(&user)
        .set(changeset)
        .execute(conn)
        .map_err(|_| CentaureissiError::RelationalDatabaseError())?;

    Ok((StatusCode::OK, ""))
}

async fn user_update_password(
    State(context): State<Arc<CentaureissiContext>>,
    Extension(user): Extension<User>,
    Json(input): Json<UserUpdatePasswordRequest>,
) -> Result<impl IntoResponse, CentaureissiError> {
    let conn = &mut context.db.get().unwrap();

    let hashed_old_password = PasswordHash::new(&user.password)?;
    let argon2 = Argon2::default();

    argon2
        .verify_password(input.old_password.as_bytes(), &hashed_old_password)
        .map_err(|_| CentaureissiError::IncorrectPasswordError())?;

    let salt = SaltString::generate(&mut OsRng);
    let argon2 = Argon2::default();
    let hashed_password = argon2
        .hash_password(input.new_password.as_bytes(), &salt)?
        .to_string();
    let changeset = UpdateUserPassword {
        password: hashed_password,
    };

    diesel::update(&user)
        .set(changeset)
        .execute(conn)
        .map_err(|_| CentaureissiError::RelationalDatabaseError())?;

    Ok((StatusCode::OK, ""))
}
