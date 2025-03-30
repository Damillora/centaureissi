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
    let dir = Maildirs::new(Path::new(&maildir));
    for mdir in dir.iter() {
        if verbose {
            println!("Checking folder {}", mdir.name);
        }
        let mut counter_size = 0;
        let mut counter = 0;
        let mut messages = Vec::<MailEntry>::new();

        if let Ok(msgs) = mdir.maildir.read() {
            for msg in msgs {
                let contents = msg.read()?;
                let file_name = msg.file_name()?;
                counter_size += contents.len();
                counter += 1;
                messages.push(MailEntry {
                    file_name: String::from(file_name),
                    contents: contents,
                });
                if verbose {
                    println!("Uploading {}", file_name);
                } 
                if counter_size >= 10_000_000 {
                    println!("Uploading {} emails", counter);
                    http_client
                        .upload_message_batch(messages)
                        .await?;
                    counter_size = 0;
                    counter = 0;
                    messages = Vec::<MailEntry>::new();
                }

                if delete_source {
                    msg.remove()?;
                }
            }
        }
    }

    Ok(())
}
