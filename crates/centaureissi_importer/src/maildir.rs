use std::path::Path;

use maildirs::Maildirs;

use crate::{
    client::{CentaureissiClient, MailEntry},
    errors::ImporterError,
};

pub async fn import_maildir(
    http_client: CentaureissiClient,
    verbose: bool,
    maildir: String,
    delete_source: bool,
) -> Result<(), ImporterError> {
    println!("Importing maildir {}", maildir);
    let mut idx = 1;

    let dir = Maildirs::new(Path::new(&maildir));
    for mdir in dir.iter() {
        if verbose {
            println!("Checking folder {}", mdir.name);
        }
        let mut counter = 0;
        let mut messages = Vec::<MailEntry>::new();

        if let Ok(msgs) = mdir.maildir.read() {
            for msg in msgs {
                let contents = msg.read()?;
                let file_name = msg.file_name()?;
                counter += contents.len();
                messages.push(MailEntry {
                    file_name: String::from(file_name),
                    contents: contents,
                });
                if verbose {
                    println!("Uploading {}", file_name);
                } else if idx % 1000 == 0 {
                    println!("Uploading message #{}", idx);
                }
                if counter >= 10_000_000 {
                    http_client
                        .upload_message_batch(messages)
                        .await?;
                    counter = 0;
                    messages = Vec::<MailEntry>::new();
                }
                idx = idx + 1;

                if delete_source {
                    msg.remove()?;
                }
            }
        }
    }

    Ok(())
}
