FROM golang:1.21-alpine

EXPOSE 8080

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o trash_bot ./cmd/bot

USER 1000:1000

ENV PORT=8080

ENTRYPOINT ["/app/trash_bot"]
