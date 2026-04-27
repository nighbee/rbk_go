package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"weather-api-t3/internal/client"
	"weather-api-t3/internal/handler"
	"weather-api-t3/internal/repository"
	"weather-api-t3/internal/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Setup DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback for local dev if not set
		dbURL = "postgres://postgres:postgres@localhost:5432/weather_db?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// 2. Setup Components
	repo := repository.NewRepository(db)
	
	httpClient := &http.Client{Timeout: 10 * time.Second}
	weatherClient := client.NewWeatherClient(httpClient)
	
	userService := service.NewUserService(repo)
	weatherService := service.NewWeatherService(repo, repo, weatherClient)
	
	h := handler.NewHandler(userService, weatherService)

	// 3. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, h.Routes()); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
