# Alterado de 1.21 para 1.23 (ou 1.22) para satisfazer o go.mod
FROM golang:1.24-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário de forma estática (ideal para containers leves)
RUN CGO_ENABLED=0 GOOS=linux go build -o mqtt-broker .

# Estágio Final (Runtime)
FROM alpine:latest

WORKDIR /root/

# Copia apenas o binário do estágio de build
COPY --from=builder /app/mqtt-broker .

# Expoe a porta padrão do MQTT
EXPOSE 1883

# Comando para rodar o mqtt-broker
CMD ["./mqtt-broker"]