use anyhow::Result;

pub async fn do_request(url: &str, method: reqwest::Method, body: Option<&str>) -> Result<String> {
    let client = reqwest::Client::new();
    let request_builder = client.request(method, url);

    let request_builder = if let Some(body) = body {
        request_builder.body(body.to_string())
    } else {
        request_builder
    };

    let response = request_builder.send().await?;
    let response_text = response.text().await?;
    Ok(response_text)
}
