package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"monitor-cripto/go-api/handlers"
	"monitor-cripto/go-api/database"
	"fmt"
)


func main() {
    database.InitDB()
    fmt.Println("Banco de dados inicializado com sucesso!")

    r := mux.NewRouter()
    r.HandleFunc("/receive-data", handlers.ReceiveCryptoData).Methods("POST")
    r.HandleFunc("/cryptos", handlers.GetCryptos).Methods("GET")  // aqui!

    log.Println("Servidor iniciado na porta 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
}
