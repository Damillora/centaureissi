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

impl CentaureissiConfig {
    pub fn get_database_url(&self) -> String {
        let mut database_url = self.data_dir.clone();
        database_url.push_str(&"/centaureissi.db");

        database_url
    }
    pub fn get_search_index_path(&self) -> String {
        let mut search_index = self.data_dir.clone();
        search_index.push_str(&"/search");

        search_index
    }
    pub fn get_blob_db_path(&self) -> String {
        let mut blob_db_file = self.data_dir.clone();
        blob_db_file.push_str(&"/blobs.db");

        blob_db_file
    }
}
