package database

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3" // driver SQLite
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./cryptos.db")
    if err != nil {
        log.Fatalf("Erro ao abrir banco de dados: %v\n", err)
    }

    // Testa a conexão
    err = DB.Ping()
    if err != nil {
        log.Fatalf("Erro ao conectar no banco de dados: %v\n", err)
    }

    log.Println("Banco de dados conectado com sucesso!")

    createTable()
}

func createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS crypto (
        id TEXT PRIMARY KEY,
        symbol TEXT,
        name TEXT,
        current_price REAL,
        market_cap_rank INTEGER,
        price_change_percentage_24h REAL
    );
    `
    _, err := DB.Exec(query)
    if err != nil {
        log.Fatalf("Erro ao criar tabela crypto: %v\n", err)
    }
    log.Println("Tabela crypto pronta.")
}
