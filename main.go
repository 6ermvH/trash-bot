package main

import (
	"os"
	"net/http"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "✅ Telegram bot is alive!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Listening on port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil)) // 🔥 блокирует main()
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})

	go InitTelegramAPI()

	server()
}
