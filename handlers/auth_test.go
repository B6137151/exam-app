package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Exam/config"
	"Exam/database"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	os.Setenv("JWT_SECRET_KEY", "mysecretkey")
	database.InitDB(config.GetConfig())
}

func TestRegisterHandler(t *testing.T) {
	database.DB.Exec("DELETE FROM users") // Clean up before running the test

	payload := `{"username":"testuser","password":"testpass"}`
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := map[string]string{"message": "User created successfully"}
	var actual map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if actual["message"] != expected["message"] {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestLoginHandler(t *testing.T) {
	// Ensure the user exists
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", "testuser", string(hashedPassword))

	payload := `{"username":"testuser","password":"testpass"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response["message"] != "Login successful" {
		t.Errorf("handler returned unexpected message: got %v want %v", response["message"], "Login successful")
	}

	if response["token"] == "" {
		t.Error("handler did not return a token")
	}
}

func TestLogoutHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogoutHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := map[string]string{"message": "Logout successful"}
	var actual map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if actual["message"] != expected["message"] {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
