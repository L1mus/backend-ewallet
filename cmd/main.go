package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load Environment Variable")
	}
	mux := http.NewServeMux()

	mux.HandleFunc("POST /hello", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello World")
		if err != nil {
			return
		}
	})

	s := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: mux,
	}

	log.Printf("Starting server %s", os.Getenv("PORT"))
	if err := s.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		log.Fatal("Server startup failed", err)
	}
}
