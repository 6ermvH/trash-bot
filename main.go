package main

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	InitTelegramAPI()
}
