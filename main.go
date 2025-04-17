package main

import (
	"os"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	InitTelegramAPI()

	go func() {
    http.ListenAndServe(":8080", nil)
	}()
}
