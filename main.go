package main

import (
	"os"
	"log"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func server() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "✅ Telegram bot is alive")
	})
	log.Printf("🌐 Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("⚠️ HTTP server error: %v", err)
	}
}

func main() {
	fmt.Println(os.Getenv("REDIS_ADDR"))
	fmt.Println(os.Getenv("REDIS_PASSWORD"))

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       1,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
			log.Fatal("Redis connect error:", err)
	}
	fmt.Println("Redis:", pong)

	go InitTelegramAPI()

	server()

	select {}
}
