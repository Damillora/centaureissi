use centaureissi_server::{config::CentaureissiConfig, http};
use config::Config;
use diesel::{prelude::*, r2d2::{self, ConnectionManager}};
use diesel_migrations::{embed_migrations, EmbeddedMigrations, MigrationHarness};
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("../../migrations");


#[tokio::main]
async fn main() {
    let config = Config::builder()
        .add_source(config::File::with_name("centaureissi"))
        .add_source(config::Environment::with_prefix("CENTAUREISSI"))
        .build()
        .unwrap()
        .try_deserialize::<CentaureissiConfig>()
        .unwrap();

    // Generate database URL from data dir
    let mut database_url = config.data_dir.clone();
    database_url.push_str(&"/centaureissi.db");

    let connection_manager = ConnectionManager::<SqliteConnection>::new(&database_url);
    let pool = r2d2::Pool::builder().build(connection_manager).expect("Failed to create connection pool");
    
    // Migrate the Database first
    let mut conn = pool.get().unwrap();
    conn.run_pending_migrations(MIGRATIONS).unwrap();

    http::serve(config, pool).await;
}
