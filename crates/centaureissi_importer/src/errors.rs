
pub enum ImporterError {
    HttpClientError(reqwest::Error),
    MaildirError(maildirs::Error),
    MaildirIoError(std::io::Error),
    UploadError(String)
}

impl From<reqwest::Error> for ImporterError {
    fn from(e: reqwest::Error) -> Self {
        Self::HttpClientError(e)
    }
}
impl From<maildirs::Error> for ImporterError {
    fn from(e: maildirs::Error) -> Self {
        Self::MaildirError(e)
    }
}
impl From<std::io::Error> for ImporterError {
    fn from(e: std::io::Error) -> Self {
        Self::MaildirIoError(e)
    }
}
impl std::fmt::Display for ImporterError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::HttpClientError(e) => write!(f, "{}", e.to_string()),
            Self::MaildirError(e) => write!(f, "{}", e.to_string()),
            Self::MaildirIoError(e) => write!(f, "{}", e.to_string()),
            Self::UploadError(e) => write!(f, "{}", e.to_string()),
        }
    }
}
