use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
// Define a structure for holding claims data used in JWT tokens
pub struct Claims {
    pub exp: usize,   // Expiry time of the token
    pub iat: usize,   // Issued at time of the token
    pub name: String, // Email associated with the token
    pub iss: String,  // Issued by
    pub sub: i32,     // User ID
    pub aud: String,  // Audience
}
