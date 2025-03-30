use std::sync::Arc;

use crate::{client::CentaureissiClient, errors::ImporterError};

pub async fn import_mbox(
    http_client: CentaureissiClient,
    verbose: bool,
    mbox: String,
) -> Result<(), ImporterError> {
    todo!()
}
