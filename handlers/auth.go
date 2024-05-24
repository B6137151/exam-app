package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"Exam/database"
	"Exam/models"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Tel      string `json:"tel"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// UserStore interface to abstract database operations
type UserStore interface {
	CreateUser(username, password, email, tel string) error
	GetUser(username string) (models.User, error)
}

// RealUserStore implements UserStore interface
type RealUserStore struct{}

func (store *RealUserStore) CreateUser(username, password, email, tel string) error {
	_, err := database.DB.Exec("INSERT INTO users (username, password, email, tel) VALUES ($1, $2, $3, $4)", username, password, email, tel)
	return err
}

func (store *RealUserStore) GetUser(username string) (models.User, error) {
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, password, email, tel FROM users WHERE username=$1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Tel)
	return user, err
}

func RegisterHandler(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		err = store.CreateUser(creds.Username, string(hashedPassword), creds.Email, creds.Tel)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		response := map[string]string{"message": "User created successfully"}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func LoginHandler(store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		storedCreds, err := store.GetUser(creds.Username)
		if err != nil || bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)) != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &Claims{
			Username: creds.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		response := map[string]string{"message": "Login successful", "token": tokenString}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})

	response := map[string]string{"message": "Logout successful"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
