use std::sync::Arc;

use diesel::{r2d2::{self, ConnectionManager}, SqliteConnection};

use crate::config::CentaureissiConfig;

#[derive(Clone)]
pub struct CentaureissiContext {
    pub config: Arc<CentaureissiConfig>,
    pub db: r2d2::Pool<ConnectionManager<SqliteConnection>>,
}