FROM golang:1.21-alpine

EXPOSE 8080

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ${CONFIG_PATH} /config/server.yaml
RUN go build -ldflags="-s -w" -o trash_bot ./cmd/bot

ENV CONFIG_PATH=./config/server.yaml

# Подключение переменных окружения
ENV REDIS_USERNAME=${REDIS_USERNAME}
ENV REDIS_PASSWORD=${REDIS_PASSWORD}
ENV TELEGRAM_APITOKEN=${TELEGRAM_APITOKEN}
ENV OPENROUTER_API_KEY=${OPENROUTER_API_KEY}

ENTRYPOINT ["/app/trash_bot"]

