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
# Kopiere das Binary und alle wichtigen Ordner (wie templates/ oder static/)
COPY --from=builder /app/server-binary .
COPY --from=builder /app/templates ./templates
# Falls du einen static-Ordner hast, nimm ihn auch mit:
# COPY --from=builder /app/static ./static

EXPOSE 7779
CMD ["./server-binary"]