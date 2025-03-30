use std::{error::Error, fs};

use persy::{Config, Persy, PersyId, ValueMode};

pub fn initialize_blobs(data_dir: String) -> Persy {
    let mut blob_db_file = data_dir;
    blob_db_file.push_str(&"/blobs.db");

    if !fs::exists(&blob_db_file).unwrap() {
        Persy::create(&blob_db_file).expect("Cannot create blob database");
    }

    let blob_db = Persy::open(blob_db_file, Config::new()).expect("Cannot open blob database");
    let blob_table_exists = blob_db
        .exists_segment(BLOB_TABLE)
        .expect("Cannot check for existence");

    if !blob_table_exists {
        {
            let mut tx = blob_db.begin().unwrap();
            tx.create_segment(BLOB_TABLE)
                .expect("Cannot create blobs segment");
            tx.create_index::<String, PersyId>(BLOB_INDEX, ValueMode::Replace)
                .expect("Cannot create index");
            let prepared = tx.prepare().unwrap();
            prepared.commit().expect("Cannot create blob structures");
        }
    }

    blob_db
}

pub const BLOB_TABLE: &str = "blobs";
pub const BLOB_INDEX: &str = "blob_index";
