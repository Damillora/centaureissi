use serde::Serialize;

#[derive(Serialize)]
pub struct LoginResponse {
    pub token: String,
}

#[derive(Serialize)]
pub struct UserProfileResponse {
    pub username: String,
}

#[derive(Serialize)]
pub struct SearchResponse {
    pub items: Vec<SearchResponseItem>,
    pub page: usize,
    pub total_pages: usize,
    pub total_items: usize,
}

#[derive(Serialize)]
pub struct SearchResponseItem {
    pub hash: String,
    pub user_id: i64,

    pub from: String,
    pub to: String,
    pub cc: String,
    pub bcc: String,
    pub subject: String,
    pub date: String,
}

#[derive(Serialize)]
pub struct StatsResponse {
    pub version: String,
    pub db_size: u64,
    pub message_count: u64,
    pub blob_db_size: u64,
    pub blob_count: u64,
    pub search_doc_count: u64,
}

#[derive(Serialize)]
pub struct MessageResponse {
    pub hash: String,
    pub date: String,

    pub from: String,
    pub to: String,
    pub cc: String,
    pub bcc: String,

    pub is_html_mail: bool,
    pub is_text_mail: bool,
    pub has_attachments: bool,
    pub subject: String,
}

#[derive(Serialize)]
pub struct MessageContentResponse {
    pub content: String,
}

#[derive(Serialize)]
pub struct MessageAttachmentResponse {
    pub items: Vec<MessageAttachmentResponseItem>,
    pub total_items: usize,
}

#[derive(Serialize)]
pub struct MessageAttachmentResponseItem {
    pub id: usize,
    pub name: String,
}
