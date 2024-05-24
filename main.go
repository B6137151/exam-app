package main

import (
	"Exam/config"
	"Exam/database"
	"Exam/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	cfg := config.GetConfig()
	database.InitDB(cfg)

	store := &handlers.RealUserStore{}

	http.HandleFunc("/register", handlers.RegisterHandler(store))
	http.HandleFunc("/login", handlers.LoginHandler(store))
	http.HandleFunc("/logout", handlers.LogoutHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
