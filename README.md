# Kafka Consumer Microservice

Este é um microserviço que atua como consumidor de mensagens do Kafka, processando payloads e encaminhando-os para outro microserviço via REST API.

## Descrição

O microserviço é responsável por:
1. Consumir mensagens de um tópico específico do Kafka
2. Processar o payload recebido (formato definido em `payload.json`)
3. Encaminhar o payload para outro microserviço via REST API (método PATCH)

## Pré-requisitos

- Go 1.24.1 ou superior
- Apache Kafka v3.6
- Docker (opcional, para ambiente de desenvolvimento)

## Estrutura do Projeto

```
.
├── README.md
├── go.mod
├── go.sum
├── main.go
├── internal/
│   ├── consumer/
│   │   └── kafka.go
│   ├── service/
│   │   └── processor.go
│   └── config/
│       └── config.go
└── .env
```

## Configuração

1. Configure as variáveis de ambiente
Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:
```env
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_USERNAME=exemplo
KAFKA_PASSWORD=exemplo
KAFKA_SASL_MECHANISM=SCRAM-SHA-256
KAFKA_SECURITY_PROTOCOL="SASL_PLAINTEXT"
KAFKA_GROUP_ID=seu-grupo
KAFKA_TOPIC=seu-topico
KAFKA_TOPIC_DLQ=seu-topico-dlq
TARGET_SERVICE_URL=http://localhost:8080/api/v1
```

## Como Executar

Avaliaremos o código com nossos sistemas

1. Crie um ambiente em docker com Kafka v3.6 para realizar seu desenvolvimento e testes.

## Dependências Principais

- github.com/confluentinc/confluent-kafka-go/v2
- github.com/spf13/viper (para gerenciamento de variáveis de ambiente)

## Formato do Payload

O payload esperado segue o formato definido em `payload.json`. O mesmo payload será encaminhado para o serviço destino via método PATCH.

## Logs

O serviço registra logs para:
- Conexão com o Kafka
- Recebimento de mensagens
- Processamento do payload
- Chamadas ao serviço destino
- Erros durante o processamento

## Contribuição

1. Faça um repositório no seu GitHub para ser compartilhado com o time.

## Guia de Execução

Siga os passos abaixo para rodar a aplicação localmente:

1. **Subir o Kafka no Docker**
   Certifique-se de que o Docker está rodando e execute:
   ```bash
   docker-compose up -d
   ```

2. **Criar o Tópico com o Script**
   Utilize o script facilitador para criar os tópicos necessários (`orders` e `orders-dlq`):
   ```bash
   ./scripts/create-topics.sh
   ```

3. **Buildar o Projeto Go e Rodar**
   Baixe as dependências e inicie a aplicação:
   ```bash
   go mod tidy
   go run main.go
   ```
   *(Ou se preferir gerar o binário: `go build -o app main.go && ./app`)*

4. **Enviar Mensagens para o Tópico**
   Em outro terminal, utilize o script de envio para testar o processamento:
   ```bash
   ./scripts/send-message.sh
   ```

### Rodando os Testes
Para executar os testes unitários:
```bash
go test ./... -v
```

