use diesel::prelude::*;

use crate::db::schema::users;

// Gets
#[derive(Queryable, Selectable)]
#[diesel(table_name = users)]
pub struct UserProfile {
    pub username: String,
}

// Gets
#[derive(Queryable, Selectable)]
#[diesel(table_name = users)]
pub struct UserLogin {
    pub username: String,
}
