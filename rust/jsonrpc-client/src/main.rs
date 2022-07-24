use jsonrpc_client_transports::RawClient;
use jsonrpc_client_transports::transports::http;
use jsonrpc_core::Params;
use serde_json::json;

#[tokio::main]
async fn main() {
    let client = http::connect::<RawClient>("http://localhost:5000").await.unwrap();
    let mut params = serde_json::Map::new();
    params.insert("key".to_string(), json!("value"));
    let result = client.call_method("initialize", Params::Map(params)).await.unwrap();
    println!("{}", result);
}
