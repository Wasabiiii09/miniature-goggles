# 1. Build-Phase
FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server-binary ./cmd/server/main.go

# 2. Runtime-Phase (schlankes Finales Image)
FROM alpine:latest
WORKDIR /app

# Kopiere das Binary und den public-Ordner mit den HTML-Dateien
COPY --from=builder /app/server-binary .
COPY --from=builder /app/public ./public

EXPOSE 7779
CMD ["./server-binary"]