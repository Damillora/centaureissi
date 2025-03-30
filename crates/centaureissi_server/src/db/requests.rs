use diesel::prelude::*;

use crate::db::schema::{messages, users};

// Creates
#[derive(Insertable)]
#[diesel(table_name = users)]
pub struct NewUser {
    pub username: String,
    pub password: String,
}

#[derive(Insertable)]
#[diesel(table_name = messages)]
pub struct NewMessage {
    pub user_id: i32,
    pub content_hash: String,
}

// Update
#[derive(AsChangeset)]
#[diesel(table_name = users)]
pub struct UpdateUser {
    pub username: String,
}
// Update
#[derive(AsChangeset)]
#[diesel(table_name = users)]
pub struct UpdateUserPassword {
    pub password: String,
}