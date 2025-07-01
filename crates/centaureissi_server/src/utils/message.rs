use mail_parser::{Message, MimeHeaders};

pub struct MessageModel {
    pub user_id: i64,

    pub from: String,
    pub to: String,
    pub cc: String,
    pub bcc: String,
    pub subject: String,
    pub timestamp_secs: i64,
    pub content: String,
    pub is_html_mail: bool,
    pub is_text_mail: bool,
    pub has_attachments: bool,
}

pub struct MessageContentModel {
    pub content: String,
}

pub struct MessageAttachModel {
    pub name: String,
    pub content: Vec<u8>,
}

pub struct MessageAttachment {
    pub id: usize,
    pub name: String,
}
pub fn create_message_model_from_message(message_user_id: i32, msg: Message) -> MessageModel {
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

    let is_html_mail = msg.html_body_count() > 0;
    let is_text_mail = msg.text_body_count() > 0;
    let has_attachments = msg.attachment_count() > 0;

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

    return MessageModel {
        user_id: i64::from(message_user_id),
        from: from_data.join(", "),
        to: to_data.join(", "),
        cc: cc_data.join(", "),
        bcc: bcc_data.join(", "),
        subject: subject_data.to_owned(),
        timestamp_secs: date_data,

        is_html_mail: is_html_mail,
        is_text_mail: is_text_mail,
        has_attachments: has_attachments,
        content: mail_contents_data.join("\n\n"),
    };
}

pub fn create_message_content_model_from_message(html: bool, msg: Message) -> MessageContentModel {
    if html {
        let mail_contents_html: Vec<String> = msg
            .html_bodies()
            .map(|item| item.text_contents().unwrap_or("").to_string())
            .collect();
        MessageContentModel {
            content: mail_contents_html.join("\n\n"),
        }
    } else {
        let mail_contents_data: Vec<String> = msg
            .text_bodies()
            .map(|item| item.text_contents().unwrap_or("").to_string())
            .collect();

        MessageContentModel {
            content: mail_contents_data.join("\n\n"),
        }
    }
}

pub fn create_message_attachment_list_from_message(msg: Message) -> Vec<MessageAttachment> {
    let mail_attachments: Vec<MessageAttachment> = msg
        .attachments()
        .enumerate()
        .map(|(i, item)| MessageAttachment {
            id: i,
            name: item.attachment_name().unwrap_or("Attachment").to_string(),
        })
        .collect();

    mail_attachments
}
pub fn create_message_attachment_content_model_from_message(
    idx: usize,
    msg: Message,
) -> MessageAttachModel {
    let mail_contents_attachment = msg.attachment(idx).unwrap();

    MessageAttachModel {
        name: mail_contents_attachment
            .attachment_name()
            .unwrap_or("Attachmennt")
            .to_string(),
        content: mail_contents_attachment.contents().to_vec(),
    }
}
