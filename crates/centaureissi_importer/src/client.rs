
use reqwest::StatusCode;

use crate::errors::ImporterError;



pub struct CentaureissiClient {
    pub server: String,
    pub token: String,
}

impl CentaureissiClient {
    pub async fn upload_message(&self, file_name: String, file: Vec<u8>) -> Result<(), ImporterError> {
        use reqwest::multipart;

        let api_route = "/api/messages/add";
        let mut api_path = self.server.clone();
        api_path.push_str(api_route);

        let eml = multipart::Part::stream(file)
            .file_name(file_name.clone());

        let form = multipart::Form::new()
            .part("file", eml);
        
        let client = reqwest::Client::new();
        let resp = client
            .post(api_path)
            .header("Authorization", format!("Bearer {}", self.token))
            .multipart(form)
            .send().await?;
        
        if resp.status() == StatusCode::OK {
            return Ok(())
        } else if resp.status() == StatusCode::BAD_REQUEST {
            let result = resp.text().await?;
            return Err(ImporterError::UploadError(format!("{}", result)))
        } else if resp.status() == StatusCode::INTERNAL_SERVER_ERROR {
            let result = resp.text().await?;
            return Err(ImporterError::UploadError(format!("{}", result)))
        } else {
            resp.error_for_status()?;
        }

        
        Ok(())
    }
}