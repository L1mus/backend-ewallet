package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/L1mus/backend-ewallet/docs"
	"github.com/L1mus/backend-ewallet/internal/config"
	"github.com/L1mus/backend-ewallet/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           E-Wallet API
// @version         1.0
// @description     This is a backend project for a web application called E-wallet.

// @contact.name   Ali Mustadji
// @contact.url    https://github.com/L1mus
// @contact.email  limustadji@gmail.com

// @license.name  MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apiKey  ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer token used for authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env. \ncause: %s", err.Error())
	}
	if err := os.MkdirAll(filepath.Join("public", "img"), os.ModePerm); err != nil {
		log.Fatalf("Failed to create upload directory: %s", err.Error())
	}
	// inisialisasi
	// gin.New()
	app := gin.Default()
	// connect ke db
	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")
	// install router
	router.InitRouter(app, db)
	// run
	// addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))
}
