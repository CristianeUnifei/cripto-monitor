package handlers

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"

    "monitor-cripto/go-api/database"
    "monitor-cripto/go-api/models"
)

func ReceiveCryptoData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Erro ao ler corpo da requisição: %v\n", err)
        http.Error(w, "Erro ao ler o corpo da requisição", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    var dados []models.CryptoData
    err = json.Unmarshal(body, &dados)
    if err != nil {
        log.Printf("Erro ao decodificar JSON: %v\n", err)
        http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
        return
    }

    tx, err := database.DB.Begin()
    if err != nil {
        log.Printf("Erro iniciando transação: %v\n", err)
        http.Error(w, "Erro iniciando transação", http.StatusInternalServerError)
        return
    }

    stmt, err := tx.Prepare(`INSERT OR REPLACE INTO crypto
        (id, symbol, name, current_price, market_cap_rank, price_change_percentage_24h)
        VALUES (?, ?, ?, ?, ?, ?)`)
    if err != nil {
        log.Printf("Erro preparando statement: %v\n", err)
        http.Error(w, "Erro preparando statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    for _, crypto := range dados {
        _, err = stmt.Exec(
            crypto.ID,
            crypto.Symbol,
            crypto.Name,
            crypto.CurrentPrice,
            crypto.MarketCapRank,
            crypto.PriceChangePercentage24h,
        )
        if err != nil {
            tx.Rollback()
            log.Printf("Erro inserindo dados no banco: %v\n", err)
            http.Error(w, "Erro inserindo dados no banco", http.StatusInternalServerError)
            return
        }

        fmt.Printf("%s (%s): $%.2f - variação 24h: %.2f%%\n",
            crypto.Name, crypto.Symbol, crypto.CurrentPrice, crypto.PriceChangePercentage24h)
    }

    err = tx.Commit()
    if err != nil {
        log.Printf("Erro salvando transação: %v\n", err)
        http.Error(w, "Erro salvando transação", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Dados recebidos e salvos com sucesso!"))
}

func GetCryptos(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT id, symbol, name, current_price, market_cap_rank, price_change_percentage_24h FROM crypto")
    if err != nil {
        log.Printf("Erro ao consultar banco: %v\n", err)
        http.Error(w, "Erro ao consultar banco", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var cryptos []models.CryptoData

    for rows.Next() {
        var c models.CryptoData
        err = rows.Scan(&c.ID, &c.Symbol, &c.Name, &c.CurrentPrice, &c.MarketCapRank, &c.PriceChangePercentage24h)
        if err != nil {
            log.Printf("Erro ao ler dados do banco: %v\n", err)
            http.Error(w, "Erro ao ler dados do banco", http.StatusInternalServerError)
            return
        }
        cryptos = append(cryptos, c)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Erro ao iterar resultados do banco: %v\n", err)
        http.Error(w, "Erro ao iterar resultados do banco", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cryptos)
}
