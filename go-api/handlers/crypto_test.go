package handlers

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestReceiveCryptoData(t *testing.T) {
    // JSON de exemplo com um array de criptos
    jsonData := `[{
        "id": "bitcoin",
        "symbol": "btc",
        "name": "Bitcoin",
        "current_price": 30000.5,
        "market_cap_rank": 1,
        "price_change_percentage_24h": 2.5
    }]`

    req, err := http.NewRequest("POST", "/receive-data", bytes.NewBuffer([]byte(jsonData)))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Cria um ResponseRecorder para capturar a resposta
    rr := httptest.NewRecorder()

    // Chama o handler
    handler := http.HandlerFunc(ReceiveCryptoData)
    handler.ServeHTTP(rr, req)

    // Verifica o status HTTP esperado
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Status do handler incorreto: esperado %v, obtido %v", http.StatusOK, status)
    }

    // Verifica se a resposta tem o texto esperado
    expected := "Dados recebidos com sucesso!"
    if rr.Body.String() != expected {
        t.Errorf("Resposta inesperada do handler: esperado %q, obtido %q", expected, rr.Body.String())
    }
}
