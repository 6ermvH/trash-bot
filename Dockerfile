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
ENV REDIS_ADDR="154.194.53.129:6379"
ENV REDIS_USERNAME="default"
ENV REDIS_PASSWORD=""

CMD ["/app/bot"]
