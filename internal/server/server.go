package server

import (
	"fmt"
	"gitlab.com/6ermvH/trash_bot/internal/config"
	"log"
	"net/http"
)

// Start runs a simple HTTP server for healthchecks
func Start(cfg config.ServerConfig) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})
	log.Printf("Health server listening on :%s", cfg.Port)
	if err := http.ListenAndServe(
		":"+cfg.Port,
		nil,
	); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
