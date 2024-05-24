package main

import (
	"Exam/config"
	"Exam/database"
	"Exam/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.GetConfig()
	database.InitDB(cfg)

	store := &handlers.RealUserStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterHandler(store))
	mux.HandleFunc("/login", handlers.LoginHandler(store))
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	// CORS middleware
	handler := cors.Default().Handler(mux)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
