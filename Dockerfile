# Etapa 1: build (imagem oficial do Go)
FROM golang:1.25.3-alpine AS builder

# Instala dependências mínimas (git, ca-certificates)
RUN apk add --no-cache git ca-certificates

# Diretório de trabalho
WORKDIR /app

# Copia os arquivos Go mod
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código
COPY . .

# Compila para binário estático (GOOS linux, GOARCH amd64)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# Etapa 2: imagem final minimalista
FROM scratch

# Copia certificados SSL e binário do build
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main

EXPOSE 8000

# Comando padrão
ENTRYPOINT ["/main"]