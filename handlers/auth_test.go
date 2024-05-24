package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Exam/models"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type MockUserStore struct {
	users map[string]models.User
}

func (store *MockUserStore) CreateUser(username, password, email, tel string) error {
	store.users[username] = models.User{Username: username, Password: password, Email: email, Tel: tel}
	return nil
}

func (store *MockUserStore) GetUser(username string) (models.User, error) {
	user, exists := store.users[username]
	if !exists {
		return models.User{}, sql.ErrNoRows
	}
	return user, nil
}

func init() {
	os.Setenv("JWT_SECRET_KEY", "mysecretkey")
}

func TestRegisterHandler(t *testing.T) {
	store := &MockUserStore{users: make(map[string]models.User)}

	payload := `{"username":"testuser","password":"testpass","email":"test@example.com","tel":"1234567890"}`
	req, err := http.NewRequest("POST", "/register", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := RegisterHandler(store)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "handler returned wrong status code")
	assert.JSONEq(t, `{"message":"User created successfully"}`, rr.Body.String(), "handler returned unexpected body")
}

func TestLoginHandler(t *testing.T) {
	store := &MockUserStore{users: make(map[string]models.User)}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
	store.CreateUser("testuser", string(hashedPassword), "test@example.com", "1234567890")

	payload := `{"username":"testuser","password":"testpass"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := LoginHandler(store)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Login successful", response["message"], "handler returned unexpected message")
	assert.NotEmpty(t, response["token"], "handler did not return a token")
}

func TestLogoutHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogoutHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")
	assert.JSONEq(t, `{"message":"Logout successful"}`, rr.Body.String(), "handler returned unexpected body")
}
