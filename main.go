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

	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
