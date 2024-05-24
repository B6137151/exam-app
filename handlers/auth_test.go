package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Exam/models"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	os.Setenv("JWT_SECRET_KEY", "mysecretkey")
}

// MockUserStore to use in tests
type MockUserStore struct {
	users map[string]string
}

func (store *MockUserStore) CreateUser(username, password string) error {
	if _, exists := store.users[username]; exists {
		return errors.New("user already exists")
	}
	store.users[username] = password
	return nil
}

func (store *MockUserStore) GetUser(username string) (models.User, error) {
	password, exists := store.users[username]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return models.User{Username: username, Password: password}, nil
}

func TestRegisterHandler(t *testing.T) {
	store := &MockUserStore{users: make(map[string]string)}

	payload := `{"username":"testuser","password":"testpass"}`
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := RegisterHandler(store)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expected := `{"message":"User created successfully"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestLoginHandler(t *testing.T) {
	store := &MockUserStore{users: make(map[string]string)}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	store.CreateUser("testuser", string(hashedPassword))

	payload := `{"username":"testuser","password":"testpass"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := LoginHandler(store)

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

	expected := `{"message":"Logout successful"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
