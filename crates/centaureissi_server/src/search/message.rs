use mail_parser::Message;
use tantivy::{doc, TantivyDocument};

pub fn create_search_document_from_message(message_id: i32, message_user_id: i32, content_hash: String, msg: Message) -> TantivyDocument {
    let from_data: Vec<String> = match msg.from() {
        Some(from) => from
            .iter()
            .map(|item| {
                if let Some(name) = item.name() {
                    format!("{} <{}>", name, item.address().unwrap())
                } else {
                    format!("{}", item.address().unwrap())
                }
            })
            .collect(),
        None => Vec::new(),
    };
    let to_data: Vec<String> = match msg.to() {
        Some(to) => to
            .iter()
            .map(|item| {
                if let Some(name) = item.name() {
                    format!("{} <{}>", name, item.address().unwrap_or(""))
                } else {
                    format!("{}", item.address().unwrap_or(""))
                }
            })
            .collect(),
        None => Vec::new(),
    };
    let cc_data: Vec<String> = match msg.cc() {
        Some(cc) => cc
            .iter()
            .map(|item| {
                if let Some(name) = item.name() {
                    format!("{} <{}>", name, item.address().unwrap())
                } else {
                    format!("{}", item.address().unwrap())
                }
            })
            .collect(),
        None => Vec::new(),
    };
    let bcc_data: Vec<String> = match msg.bcc() {
        Some(bcc) => bcc
            .iter()
            .map(|item| {
                if let Some(name) = item.name() {
                    format!("{} <{}>", name, item.address().unwrap())
                } else {
                    format!("{}", item.address().unwrap())
                }
            })
            .collect(),
        None => Vec::new(),
    };
    let subject_data = msg.subject().unwrap_or("");
    let date_data = msg
        .date()
        .unwrap_or(&mail_parser::DateTime::from_timestamp(0))
        .to_timestamp();
    let mail_contents_data: Vec<String> = msg
        .text_bodies()
        .map(|item| item.text_contents().unwrap_or("").to_string())
        .collect();

    let schema = crate::search::get_schema();

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
        from => from_data.join(", "),
        to => to_data.join(", "),
        cc => cc_data.join(", "),
        bcc => bcc_data.join(", "),
        subject => subject_data,
        date => tantivy::DateTime::from_timestamp_secs(date_data),
        content => mail_contents_data.join("\n\n"),
    );
}