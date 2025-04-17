FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bot

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates && apt-get clean

WORKDIR /app
COPY --from=builder /app/bot .

ENV TELEGRAM_APITOKEN=""
ENV REDIS_ADDR="redis:6379"

CMD ["./bot"]

