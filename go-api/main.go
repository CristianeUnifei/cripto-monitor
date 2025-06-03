package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"monitor-cripto/go-api/handlers"
	"monitor-cripto/go-api/database"
)

type CryptoData struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	MarketCapRank            int     `json:"market_cap_rank"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
}

func fetchCryptoPage(page int, perPage int) ([]CryptoData, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=%d&page=%d&sparkline=false", perPage, page)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []CryptoData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func fetchAllCryptos() ([]CryptoData, error) {
	var allData []CryptoData
	pages := []int{1, 2, 3}     // páginas que quer buscar
	perPage := 100              // máximo permitido pela API (até 250 é possível via paginação)
	for _, p := range pages {
		data, err := fetchCryptoPage(p, perPage)
		if err != nil {
			return nil, err
		}
		allData = append(allData, data...)
	}
	return allData, nil
}

func postCryptoData(data []CryptoData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/receive-data", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao enviar dados: %s", resp.Status)
	}
	return nil
}

func startAutoUpdate(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			log.Println("Buscando dados atualizados da CoinGecko (múltiplas criptos)...")
			data, err := fetchAllCryptos()
			if err != nil {
				log.Println("Erro ao buscar dados:", err)
				continue
			}

			err = postCryptoData(data)
			if err != nil {
				log.Println("Erro ao enviar dados para /receive-data:", err)
				continue
			}

			log.Println("Dados atualizados com sucesso!")
		}
	}()
}

func main() {
	database.InitDB()
	fmt.Println("Banco de dados inicializado com sucesso!")

	r := mux.NewRouter()
	r.HandleFunc("/receive-data", handlers.ReceiveCryptoData).Methods("POST")
	r.HandleFunc("/cryptos", handlers.GetCryptos).Methods("GET")

	// Atualiza a cada 10 minutos (ajuste conforme quiser)
	startAutoUpdate(10 * time.Minute)

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
