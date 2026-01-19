FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /trash-bot ./cmd

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /trash-bot .
COPY config/base.yaml ./config/base.yaml

EXPOSE 8080

CMD ["./trash-bot"]
