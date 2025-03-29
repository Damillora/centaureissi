use diesel::prelude::*;

use crate::db::schema::users;

// Creates
#[derive(Insertable)]
#[diesel(table_name = users)]
pub struct NewUser {
    pub username: String,
    pub password: String,
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