use crate::{blobs::BLOB_TABLE, config::CentaureissiConfig};

pub fn compress_payloads(config: CentaureissiConfig, blob_db: persy::Persy) {
    for (read_id, content) in blob_db.scan(BLOB_TABLE).unwrap() {
        let mut uncompressed = Vec::<u8>::new();

        let is_compressed = zstd::stream::copy_decode(&*content, &mut uncompressed).is_err();

        if !is_compressed {
            if config.verbose {
                println!("Found uncompressed data: {}", read_id);
            }
            // ZSTD compression
            let mut compressed_data = Vec::<u8>::new();
            zstd::stream::copy_encode(&*content, &mut compressed_data, 0).unwrap();

            // Insert email to blob transaction
            let mut blob_write_txn = blob_db.begin().unwrap();
            blob_write_txn
                .update(BLOB_TABLE, &read_id, &compressed_data)
                .unwrap();
    
            let prepared = blob_write_txn.prepare().unwrap();
            prepared.commit().unwrap();
        } else {
            if config.verbose {
                println!("Found compressed data: {}", read_id);
            }
        }
    }
}
