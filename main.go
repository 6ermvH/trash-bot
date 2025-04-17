package main

import (
	"os"
	"net/http"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func server() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // ✅ Явно указываем 200 OK
		fmt.Fprintln(w, "✅ Telegram bot is alive!")
	})

	log.Printf("🌐 Listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("❌ HTTP server error: %v", err)
	}
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	server()

	for true {
		InitTelegramAPI()
	}
}
