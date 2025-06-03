package models

type CryptoData struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCapRank            int     `json:"market_cap_rank"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}
