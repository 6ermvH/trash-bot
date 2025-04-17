FROM golang:1.21 as builder

WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/bot .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    redis-server \
    && apt-get clean

WORKDIR /app

COPY --from=builder /app/bot /app/bot

RUN chmod +x /app/bot

ENV TELEGRAM_APITOKEN=""
ENV REDIS_ADDR="redis:6379"

CMD ["/app/bot"]
