FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot

# Финальный образ с Redis + ботом
FROM debian:bookworm-slim

# Устанавливаем Redis и сертификаты для HTTPS
RUN apt-get update && apt-get install -y \
    redis-server \
    ca-certificates \
    && apt-get clean

WORKDIR /app
COPY --from=builder /app/bot .

# Устанавливаем переменные окружения
ENV REDIS_ADDR=localhost:6379
ENV PORT=8080

# Запускаем Redis и бота
CMD bash -c "redis-server --daemonize yes && ./bot"

