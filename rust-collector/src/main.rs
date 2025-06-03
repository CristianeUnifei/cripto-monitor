use reqwest::Client;
use serde::Deserialize;
use serde::Serialize;

#[derive(Debug, Deserialize, Serialize)]
struct CryptoData {
    id: String,
    symbol: String,
    name: String,
    current_price: f64,
    market_cap_rank: u32,
    price_change_percentage_24h: f64,
}

#[derive(Serialize)]
struct MeuTipo {
    campo1: String,
    campo2: i32,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let url = "https://api.coingecko.com/api/v3/coins/markets";
    let params = [
        ("vs_currency", "usd"),
        ("order", "market_cap_desc"),
        ("per_page", "10"),
        ("page", "1"),
        ("sparkline", "false"),
        ("price_change_percentage", "24h"),
    ];

    let client = Client::new();
    let response = client
        .get(url)
        .query(&params)
        .header("User-Agent", "monitor-cripto/0.1") // <- header importante
        .send()
        .await?;

    let status = response.status();
    let text = response.text().await?;

    println!("STATUS: {}", status);
    println!("BODY:\n{}", text);

    if status.is_success() {
        let cryptos: Vec<CryptoData> = serde_json::from_str(&text)?;

        for crypto in &cryptos {
            println!(
                "{} ({}) - ${} - rank {} - variação 24h: {}%",
                crypto.name,
                crypto.symbol,
                crypto.current_price,
                crypto.market_cap_rank,
                crypto.price_change_percentage_24h
            );
        }

        // Enviar para API Go
        let post_url = "http://localhost:8080/receive-data"; // Altere se necessário
        let res = client
            .post(post_url)
            .json(&cryptos)
            .send()
            .await?;

        println!("Resposta do backend: {:?}", res.status());
    } else {
        println!("Erro: status da resposta HTTP foi {}", status);
    }

    Ok(())
}
