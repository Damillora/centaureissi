use centaureissi_importer::{client::CentaureissiClient, imap::import_imap, maildir::import_maildir, mbox::import_mbox};
use clap::{ArgGroup, Parser};

/// Imports emails into the centaureissi system
#[derive(Parser, Debug)]
#[command(version, about, long_about = None)]
#[clap(group(
            ArgGroup::new("import")
                .required(true)
                .args(&["mbox", "maildir", "imap_server"]),
        ))]
struct ImporterConfig {
    /// Verbose logging
    #[arg(short, long, default_value_t = false)]
    verbose: bool,

    /// URL of centaureissi server
    #[arg(short, long)]
    server: String,

    /// Authentication token for server
    #[arg(short, long)]
    token: String,

    /// Delete source email after importing
    #[arg(long, default_value_t = false)]
    delete_source: bool,

    /// Import mbox file
    #[arg(short, long)]
    mbox: Option<String>,

    /// Import maildir
    #[arg(short('d'), long)]
    maildir: Option<String>,

    /// IMAP Server
    #[arg(short, long, group = "imap")]
    imap_server: Option<String>,

    /// IMAP Port
    #[arg(long, group = "imap", default_value_t = 143)]
    imap_port: u16,

    /// IMAP Username
    #[arg(short('u'), long, group = "imap")]
    imap_username: Option<String>,

    #[arg(short('p'), long, group = "imap")]
    imap_password: Option<String>

}

#[tokio::main]
async fn main() {
    let config = ImporterConfig::parse();
    let client = CentaureissiClient {
            server: config.server.clone(),
            token: config.token.clone(),
        };
    
    if config.mbox.is_some() {
        let mbox_file = config.mbox.unwrap();
        let result = import_mbox(client, config.verbose, mbox_file).await;
        match result {
            Ok(_) => (),
            Err(e) => println!("Error occured on import: {}", e.to_string())
        }
    } else if config.maildir.is_some() {
        let maildir = config.maildir.unwrap();
        let result = import_maildir(client, config.verbose, maildir, config.delete_source).await;
        match result {
            Ok(_) => (),
            Err(e) => println!("Error occured on import: {}", e.to_string())
        }
    } else if config.imap_server.is_some() {
        let imap_server = config.imap_server.unwrap();
        let imap_port = config.imap_port;
        let imap_username = config.imap_username.unwrap();
        let imap_password = config.imap_password.unwrap();
        let result = import_imap(client, config.verbose, imap_server, imap_port, imap_username, imap_password).await;
        match result {
            Ok(_) => (),
            Err(e) => println!("Error occured on import: {}", e.to_string())
        }
    }
}
