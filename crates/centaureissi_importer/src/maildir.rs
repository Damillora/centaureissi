use std::path::Path;

use maildirs::{MaildirEntry, Maildirs, MaildirsEntry};
use tokio::fs;

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

    let mut counter_size = 0;
    let mut counter = 0;
    let mut messages = Vec::<MailEntry>::new();
    let mut paths = Vec::<String>::new();

    let msgs = dir.iter().map(|s| s.maildir).flat_map(|dir| {
        let mdir = dir.read().unwrap();

        mdir.collect::<Vec<MaildirEntry>>()
    });

    let mut msgs_iter = msgs.peekable();

    while let Some(msg) = msgs_iter.next() {
        let contents = msg.read()?;
        let file_name = msg.file_name()?;
        let path = msg.path().to_owned().clone();
        counter_size += contents.len();
        counter += 1;
        messages.push(MailEntry {
            file_name: String::from(file_name),
            contents: contents,
        });
        paths.push(String::from(path.to_str().unwrap()));
        if verbose {
            println!("Uploading {}", file_name);
        }
        if counter_size >= 10_000_000 || counter >= 100 || msgs_iter.peek().is_none() {
            println!("Uploading {} emails", counter);
            http_client.upload_message_batch(messages).await?;
            counter_size = 0;
            counter = 0;
            messages = Vec::<MailEntry>::new();

            if delete_source {
                for path in paths {
                    fs::remove_file(path).await?;
                }
            }

            paths = Vec::<String>::new();
        }
    }

    Ok(())
}
