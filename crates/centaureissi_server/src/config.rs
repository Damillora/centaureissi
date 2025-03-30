use serde::{Deserialize, Serialize};

/// Configuration []
#[derive(Serialize, Deserialize)]
pub struct CentaureissiConfig {
    /// Data directory to store database and blobs
    pub data_dir: String, 
    /// Auth secret used to sign JWTs
    pub auth_secret: String,
    /// Disable registration on the instance
    pub disable_registration: bool,
    /// Be verbose on logging
    pub verbose: bool,
}