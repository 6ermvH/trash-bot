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
	go InitTelegramAPI()

	startHTTPHealthServer()
}

func startHTTPHealthServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bot is alive")
	})

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
