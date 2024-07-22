package main

// client responsável por fazer a requisição para o servidor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

//

const (
	url        = "http://localhost:8080/cotacao"
	timeoutReq = 300 * time.Millisecond
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutReq)
	defer cancel()

	//criar a requisição
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
	}

	//enviar a requisição
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		//verificar se o erro é de timeout
		if err == context.DeadlineExceeded {
			log.Fatalf("Tempo de resposta excedido: %v", err)
		}
		log.Fatalf("Erro ao enviar a requisição: %v", err)
	}
	defer res.Body.Close()

	//verifica a resposta
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Erro ao enviar a requisição: %v", res.Status)
	}

	//Lendo a resposta
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}

	//fmt.Println(string(body))

	//Salvando a cotacao na struct
	var cotacao map[string]string
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Fatalf("Erro ao decodificar o corpo da resposta: %v", err)
	}

	//fmt.Println("Cotação:", cotacao)

	//Abrir o arquivo em modo edicao
	_, err = os.Stat("cotacao.txt")
	fileExists := err == nil

	file, err := os.OpenFile("cotacao.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	if fileExists {
		_, err = file.WriteString(fmt.Sprintf("%s - Dólar: %s\n", currentTime, cotacao["Dólar"]))
		if err != nil {
			log.Fatalf("Erro ao escrever no arquivo: %v", err)
		}
	} else {
		_, err = file.WriteString(fmt.Sprintf("%s - Dólar: %s\n", currentTime, cotacao["Dólar"]))
		if err != nil {
			log.Fatalf("Erro ao escrever no arquivo: %v", err)
		}
	}
}
