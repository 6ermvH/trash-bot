FROM golang:1.21-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . ./

# Собираем приложение
RUN go build -ldflags="-s -w" -o trash_bot ./cmd/bot

ENV PORT=8020

# Точка входа
ENTRYPOINT ["/app/trash_bot"]

