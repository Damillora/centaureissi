use std::ops::RangeFull;

use diesel::{
    r2d2::{ConnectionManager, Pool}, ExpressionMethods, QueryDsl, RunQueryDsl, SelectableHelper, SqliteConnection
};
use persy::{IndexIter, PersyId};

use crate::{
    blobs::BLOB_INDEX,
    config::CentaureissiConfig, db::{models::{Messages, User}, requests::NewMessage},
};

pub fn rebuild_messages(
    config: CentaureissiConfig,
    default_username: String,
    blob_db: persy::Persy,
    rdb: Pool<ConnectionManager<SqliteConnection>>,
) {
    let conn = &mut rdb.get().unwrap();
    let default_user_id: i32;
    {
        use crate::db::schema::users::dsl::*;

        let user = users
            .filter(username.eq(default_username))
            .select(User::as_select())
            .first(conn)
            .unwrap();

        default_user_id = user.id;
    }

    if default_user_id == 0 {
        panic!("Default username not found!");
    }

    let iter: IndexIter<String, PersyId> = blob_db.range(BLOB_INDEX, RangeFull::default()).unwrap();
    for (hash, _) in iter {
        use crate::db::schema::messages::dsl::*;

        let message_count: i64 = messages
            .filter(content_hash.eq(&hash))
            .count()
            .get_result(conn)
            .unwrap();

        if message_count == 0 {
                // Insert into database
            let new_message = NewMessage {
                user_id: default_user_id,
                content_hash: hash,
            };

            let message = diesel::insert_into(messages)
                .values(&new_message)
                .returning(Messages::as_returning())
                .get_result(conn)
                .unwrap();

            if config.verbose {
                println!("Recovered missing message: {}", message.content_hash);
            }
        }
    }
}
