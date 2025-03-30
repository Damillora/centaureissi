use std::sync::RwLock;

use diesel::{
    SqliteConnection,
    r2d2::{self, ConnectionManager},
};
use tantivy::{Index, IndexReader, IndexWriter};

use crate::config::CentaureissiConfig;

pub struct CentaureissiContext {
    pub config: CentaureissiConfig,
    pub db: r2d2::Pool<ConnectionManager<SqliteConnection>>,
    pub search_index: Index,
    pub search_writer: RwLock<IndexWriter>,
    pub search_reader: IndexReader,
    pub blob_db: persy::Persy,
}
