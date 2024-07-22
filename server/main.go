package main

// servidor

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DollarExchange struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

const (
	url                = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	timeoutMax_Request = 200 * time.Millisecond
	timeoutMax_Persist = 10 * time.Millisecond
)

func main() {
	http.HandleFunc("/cotacao", handleCotacaoDolar)
	http.ListenAndServe(":8080", nil)
}

// função para buscar os dados do dolar do site da API
func handleCotacaoDolar(w http.ResponseWriter, r *http.Request) {
	//Criando o contexto usando o timeout para aguardar a resposta da requisição
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMax_Request)
	defer cancel()

	//Criando a requisição para o site da API
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Fprintf(w, "Erro ao criar a requisição: %v", err)
		return
	}

	// Fazendo a requisição e tratando o erro
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//verificar se o erro é de timeout
		if err == context.DeadlineExceeded {
			log.Fatalf("Tempo de resposta excedido: %v", err)
		}
		fmt.Fprintf(w, "Erro ao fazer a requisição: %v", err)
		return
	}
	defer res.Body.Close()

	// Lendo o body da resposta com a decodificação do JSON
	var exchange DollarExchange
	err = json.NewDecoder(res.Body).Decode(&exchange)
	if err != nil {
		fmt.Fprintf(w, "Erro ao decodificar o JSON: %v", err)
		return
	}

	//Criando o contexto usando o timeout para salvar a contacao no banco de dados
	ctx_DB, cancel_DB := context.WithTimeout(context.Background(), timeoutMax_Persist)
	defer cancel_DB()

	err = SaveCotacao(ctx_DB, exchange.Usdbrl.Bid)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Printf("Erro: Timeout ao salvar a cotação no banco de dados: %v", err)
		} else {
			log.Printf("Erro ao salvar a cotação no banco de dados: %v", err)
		}
		http.Error(w, fmt.Sprintf("Erro ao salvar a cotação: %v", err), http.StatusInternalServerError)
		return
	}

	// Retornando o JSON para o client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"Dólar": exchange.Usdbrl.Bid}
	json.NewEncoder(w).Encode(response)

}

// Funcao que salva os dados no sqlite
func SaveCotacao(ctx context.Context, cotacao string) error {
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		return err
	}
	defer db.Close()

	// Criando a tabela
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS cotacao (id INTEGER PRIMARY KEY AUTOINCREMENT, cotacao TEXT, data DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		return err
	}

	// Inserindo os dados
	_, err = db.ExecContext(ctx, "INSERT INTO cotacao (cotacao) VALUES (?)", cotacao)
	if err != nil {
		return err
	}
	return nil
}
