# Estágio de construção
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependência e baixar
COPY go.mod go.sum ./
RUN go mod download

# Copiar o resto do código
COPY . .

# Construir a aplicação
RUN go build -o main ./cmd/api/main.go

# Estágio final (imagem leve)
FROM alpine:latest

WORKDIR /app

# Instalar dependências básicas
RUN apk --no-cache add ca-certificates

# Copiar o binário do estágio de builder
COPY --from=builder /app/main .

# Expor a porta 8080 (Gin default no main.go)
EXPOSE 8080

# Rodar a aplicação
CMD ["./main"]
