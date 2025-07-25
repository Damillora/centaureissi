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
pub struct SignPathRequest {
    pub path: String,
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

#[derive(Deserialize)]
pub struct SearchRequest {
    pub q: String,
    pub page: Option<usize>,
    pub per_page: Option<usize>,
}

#[derive(Deserialize)]
pub struct DownloadTokenRequest {
    pub token: String,
}
