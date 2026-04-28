# Estágio de Compilação (Build)
FROM golang:1.21-alpine AS builder

# Define o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copia o código fonte
COPY . .

# Compila o binário de forma estática (ideal para containers leves)
RUN CGO_ENABLED=0 GOOS=linux go build -o broker .

# Estágio Final (Runtime)
FROM alpine:latest

WORKDIR /root/

# Copia apenas o binário do estágio de build
COPY --from=builder /app/broker .

# Expoe a porta padrão do MQTT
EXPOSE 1883

# Comando para rodar o broker
CMD ["./broker"]