use mail_parser::Message;
use tantivy::{TantivyDocument, doc};

use crate::utils::message::create_message_model_from_message;

pub fn create_search_document_from_message(
    message_user_id: i32,
    content_hash: String,
    msg: Message,
) -> TantivyDocument {
    let schema = crate::search::get_schema();

    let message_model = create_message_model_from_message(message_user_id, msg);

    // Schema Fields
    let hash = schema.get_field("hash").unwrap();
    let user_id = schema.get_field("user_id").unwrap();
    let from = schema.get_field("from").unwrap();
    let to = schema.get_field("to").unwrap();
    let cc = schema.get_field("cc").unwrap();
    let bcc = schema.get_field("bcc").unwrap();
    let subject = schema.get_field("subject").unwrap();
    let date = schema.get_field("date").unwrap();
    let content = schema.get_field("content").unwrap();

    return doc!(
        hash => content_hash,
        user_id => i64::from(message_user_id),
        from => message_model.from,
        to => message_model.to,
        cc => message_model.cc,
        bcc => message_model.bcc,
        subject => message_model.subject,
        date => tantivy::DateTime::from_timestamp_secs(message_model.timestamp_secs),
        content => message_model.content,
    );
}
