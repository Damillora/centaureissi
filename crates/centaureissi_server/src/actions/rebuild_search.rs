use diesel::{
    QueryDsl, RunQueryDsl, SelectableHelper, SqliteConnection,
    r2d2::{ConnectionManager, Pool},
};
use mail_parser::MessageParser;
use persy::PersyId;

use crate::{
    blobs::{BLOB_INDEX, BLOB_TABLE},
    config::CentaureissiConfig,
    db::models::Messages,
    search::{self, initialize_search},
};

pub fn rebuild_search_index(
    config: CentaureissiConfig,
    blob_db: persy::Persy,
    rdb: Pool<ConnectionManager<SqliteConnection>>,
    search_writer: &mut tantivy::IndexWriter,
) {
    let search = initialize_search(&config);

    let conn = &mut rdb.get().unwrap();

    {
        use crate::db::schema::messages::dsl::*;

        let all_messages = messages
            .select(Messages::as_select())
            .get_results(conn)
            .unwrap();

        let mut counter: usize = 1;

        for message_item in all_messages {
            let blob_id = blob_db
                .one::<String, PersyId>(BLOB_INDEX, &message_item.content_hash)
                .unwrap()
                .unwrap();

            let contents = blob_db.read(BLOB_TABLE, &blob_id).unwrap().unwrap();

            let mut uncompressed = Vec::<u8>::new();

            let is_compressed = !zstd::stream::copy_decode(&*contents, &mut uncompressed).is_err();

            if !is_compressed {
                uncompressed = contents;
            }

            let parsed_msg = MessageParser::default().parse(&uncompressed);
            if let Some(msg) = parsed_msg {
                let search_doc = search::message::create_search_document_from_message(
                    message_item.user_id,
                    message_item.content_hash.clone(),
                    msg,
                );

                search_writer.add_document(search_doc).unwrap();
                search_writer.commit().unwrap();

                if config.verbose {
                    println!("Reindexed message: {}", &message_item.content_hash);
                }
            } else {
                panic!("Cannot parse message: {}", &message_item.content_hash);
            }

            if counter % 100 == 0 && !config.verbose {
                println!("Indexed message #{}", counter)
            }
            counter = counter + 1;
        }
    }
}
