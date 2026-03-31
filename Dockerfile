# ====== Stage 1: Build do frontend (Vite) ======
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copia apenas o package.json/package-lock para cache de dependências
COPY frontend/package*.json ./
RUN npm install

# Copia o restante do frontend e builda
COPY frontend/ ./
RUN npm run build

# ====== Stage 2: Build do backend (Go) ======
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Dependências nativas minimalistas
RUN apk add --no-cache git ca-certificates

# Copia TODO o código do backend de uma vez
COPY backend/ ./backend/

WORKDIR /app/backend

# Baixa dependências definidas no go.mod (no seu caso, só stdlib)
RUN go mod download

# Compila o binário (nome: server)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/server

# ====== Stage 3: Imagem final ======
FROM alpine:3.19

WORKDIR /app

# Certificados para HTTPS
RUN apk add --no-cache ca-certificates

# Copia binário e dist do frontend
COPY --from=backend-builder /app/server ./server
COPY --from=frontend-builder /app/frontend/dist ./dist

# Variáveis de ambiente padrão
ENV PORT=8080
ENV FRONTEND_DIST_DIR=/app/dist

EXPOSE 8080

CMD ["./server"]
