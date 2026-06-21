package app

import (
	"ShortLink/internal/handlers"
	"log"
	"net/http"
)

func Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthHandler)
	log.Println("The server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server run error: %v", err)
	}
}