use std::fs;

use tantivy::{
    Index, TantivyError,
    schema::{INDEXED, STORED, STRING, Schema, TEXT},
};

pub fn get_schema() -> Schema {
    let mut schema_builder = Schema::builder();

    schema_builder.add_i64_field("id", INDEXED | STORED);
    schema_builder.add_text_field("hash", STRING | STORED);
    schema_builder.add_i64_field("user_id", INDEXED | STORED);

    schema_builder.add_text_field("from", STRING | STORED);
    schema_builder.add_text_field("to", STRING | STORED);
    schema_builder.add_text_field("cc", STRING | STORED);
    schema_builder.add_text_field("bcc", STRING | STORED);
    schema_builder.add_text_field("subject", TEXT | STORED);
    schema_builder.add_date_field("date", INDEXED | STORED);
    schema_builder.add_text_field("content", TEXT | STORED);
    return schema_builder.build();
}
pub fn initialize_search(data_dir: String) -> Index {
    let mut search_index = data_dir;
    search_index.push_str(&"/search");

    if !fs::exists(&search_index).unwrap() {
        fs::create_dir(&search_index).expect("Cannot create index folder!");
    }

    let schema = get_schema();

    // If index already exists, open the existing one
    let index = Index::create_in_dir(&search_index, schema.clone())
        .or_else(|error| match error {
            TantivyError::IndexAlreadyExists => Ok(Index::open_in_dir(&search_index)?),
            _ => Err(error),
        })
        .expect("Cannot create index!");

    index
}
