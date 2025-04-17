FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot

FROM debian:bookworm-slim

# Установка Redis + корневых сертификатов
RUN apt-get update && apt-get install -y \
    redis-server \
    ca-certificates \
    && apt-get clean

WORKDIR /app

COPY --from=builder /app/bot .

ENV TELEGRAM_APITOKEN="6317398679:AAE5pVghUpRGGagOsxebvlT3IqTOmcWXaxA"
ENV REDIS_ADDR=localhost:6379
ENV PORT=8080

# 👇 запускаем Redis и Telegram-бота
CMD bash -c "redis-server --daemonize yes && ./bot"

