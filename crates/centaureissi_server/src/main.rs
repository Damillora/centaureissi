use centaureissi_server::{actions, blobs, config::CentaureissiConfig, http, search::initialize_search};
use clap::{Parser, Subcommand};
use config::Config;
use diesel::{
    prelude::*,
    r2d2::{self, ConnectionManager},
};
use diesel_migrations::{EmbeddedMigrations, MigrationHarness, embed_migrations};
use tantivy::ReloadPolicy;
pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("../../migrations");

#[derive(Parser)]
#[command(version, about, long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    Serve,
    Compress,
}
#[tokio::main]
async fn main() {
    let cli = Cli::parse();

    let config = Config::builder()
        .add_source(config::File::with_name("centaureissi"))
        .add_source(config::Environment::with_prefix("CENTAUREISSI"))
        .build()
        .unwrap()
        .try_deserialize::<CentaureissiConfig>()
        .unwrap();

    let database_url = config.get_database_url();

    let connection_manager = ConnectionManager::<SqliteConnection>::new(&database_url);
    let pool = r2d2::Pool::builder()
        .build(connection_manager)
        .expect("Failed to create connection pool");

    // Migrate the Database first
    let mut conn = pool.get().unwrap();
    conn.run_pending_migrations(MIGRATIONS).unwrap();

    // Search
    let search = initialize_search(&config);
    let index_writer = search.writer(50_000_000).unwrap();
    let index_reader = search
        .reader_builder()
        .reload_policy(ReloadPolicy::OnCommitWithDelay)
        .try_into()
        .unwrap();

    // Blob Storage
    let blob_db = blobs::initialize_blobs(&config);

    match &cli.command {
        Some(Commands::Compress) => actions::compress::compress_payloads(config, blob_db),
        Some(Commands::Serve) => http::serve(config, pool, search, index_writer, index_reader, blob_db).await,
        None => http::serve(config, pool, search, index_writer, index_reader, blob_db).await,
    }
}
