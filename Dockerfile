FROM golang:1.21-alpine

EXPOSE 8080

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o trash_bot ./cmd/bot

# Подключение переменных окружения
ENV PORT=8080
ENV REDIS_ADDR=${REDIS_ADDR}
ENV REDIS_USERNAME=${REDIS_USERNAME}
ENV REDIS_PASSWORD=${REDIS_PASSWORD}
ENV TELEGRAM_APITOKEN=${TELEGRAM_APITOKEN}
ENV OPENROUTER_API_KEY=${OPENROUTER_API_KEY}

ENTRYPOINT ["/app/trash_bot"]

