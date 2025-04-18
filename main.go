package main

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	InitTelegramAPI()
}
