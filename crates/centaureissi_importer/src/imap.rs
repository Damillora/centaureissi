
use crate::{client::CentaureissiClient, errors::ImporterError};

pub async fn import_imap(
    http_client: CentaureissiClient,
    verbose: bool,
    imap_server: String,
    imap_port: u16,
    imap_username: String,
    imap_password: String,
) -> Result<(), ImporterError> {
    todo!()
}
