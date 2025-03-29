use chrono::NaiveDateTime;
use diesel::prelude::*;

use crate::db::schema::{messages, user_tokens, users};

/// Model for User
#[derive(Clone, Queryable, Identifiable, Selectable)]
#[diesel(table_name = users)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct User {
    pub id: i32,
    pub username: String,
    pub password: String,
}

#[derive(Clone, Queryable, Identifiable, Selectable)]
#[diesel(table_name = user_tokens)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct UserTokens {
    pub id: i32,
    pub user_id: i32,
    pub token: String,

    pub revoked_at: Option<NaiveDateTime>,
}

#[derive(Clone, Queryable, Identifiable, Selectable)]
#[diesel(table_name = messages)]
#[diesel(check_for_backend(diesel::sqlite::Sqlite))]
pub struct Messages {
    pub id: i32,
    pub user_id: i32,
    pub content_hash: String,
}

