package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/L1mus/backend-ewallet/internal/infrastructure/config"
	"github.com/L1mus/backend-ewallet/internal/infrastructure/database/postgres"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load Environment Variable")
	}

	// DB configuration
	configuration, err := config.LoadConfiguration()
	if err != nil {
		log.Printf("Failed to load configuration database")
	}

	// DB connection
	db, err := postgres.DBConnection(configuration)
	if err != nil {
		log.Fatalf("Connection error : %s", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database")
		}
	}(db)
	log.Println("Database Connected")

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
