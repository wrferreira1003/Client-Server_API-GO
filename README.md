### Desafio de Implementação de Webserver HTTP, Contextos, Banco de Dados e Manipulação de Arquivos com Go

Neste projeto, foram desenvolvidos dois sistemas em Go: client.go e server.go, aplicando conceitos avançados de webserver HTTP, contextos, banco de dados SQLite e manipulação de arquivos.

#### Estrutura do Projeto

##### Arquivos

  - client.go: Responsável por realizar uma requisição HTTP ao servidor solicitando a cotação do dólar.

  - server.go: Responsável por consumir a API de câmbio de Dólar e Real, retornando o resultado em formato JSON para o cliente, e registrando cada cotação recebida no banco de dados SQLite.

#### Requisitos Implementados

###### 1. Requisição HTTP
  O client.go faz uma requisição HTTP ao server.go para obter a cotação do dólar. A resposta é o valor atual do câmbio (campo "bid" do JSON).

###### 2. Consumo de API Externa
  O server.go consome a API de câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL. O resultado é enviado de volta ao client.go em formato JSON.

###### 3. Uso de Contextos
  
  - O server.go utiliza o package "context" para registrar no banco de dados SQLite cada cotação recebida.
  - Timeout máximo para chamar a API de cotação do dólar: 200ms
  - Timeout máximo para persistir os dados no banco: 10ms

  - O client.go utiliza o package "context" para receber o resultado do server.go.
  - Timeout máximo para receber o resultado: 300ms
  - Caso os tempos de execução excedam os limites estabelecidos, erros são retornados nos logs.

###### 4. Persistência em Banco de Dados
    
  - O server.go registra cada cotação recebida no banco de dados SQLite.

###### 5. Manipulação de Arquivos
  
  - O client.go salva a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}.

###### 6. Endpoint

  - O endpoint gerado pelo server.go é /cotacao e a porta utilizada pelo servidor HTTP é a 8080.
