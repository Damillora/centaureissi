use serde::Deserialize;

#[derive(Deserialize)]
pub struct NewUserRequest {
    pub username: String,
    pub password: String,
}

#[derive(Deserialize)]
pub struct LoginRequest {
    pub username: String,
    pub password: String,
}

#[derive(Deserialize)]
pub struct UserUpdateRequest {
    pub username: String,
}

#[derive(Deserialize)]
pub struct UserUpdatePasswordRequest {
    pub old_password: String,
    pub new_password: String,
}