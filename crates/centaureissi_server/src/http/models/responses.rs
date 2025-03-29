use serde::Serialize;


#[derive(Serialize)]
pub struct LoginResponse{
    pub token: String,
}

#[derive(Serialize)]
pub struct UserProfileResponse {
    pub username: String,
}