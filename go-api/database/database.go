package database

import (
    "database/sql"
    "log"

     _ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
    var err error
    DB, err = sql.Open("sqlite3", "./cryptos.db")
    if err != nil {
        log.Fatal(err)
    }

    createTableSQL := `
    CREATE TABLE IF NOT EXISTS crypto (
        id TEXT PRIMARY KEY,
        symbol TEXT,
        name TEXT,
        current_price REAL,
        market_cap_rank INTEGER,
        price_change_percentage_24h REAL
    );`

    _, err = DB.Exec(createTableSQL)
    if err != nil {
        log.Fatalf("Erro ao criar tabela: %v", err)
    }
}
