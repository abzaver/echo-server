package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateUser(t *testing.T) {
	// Setup the handler
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")

	// Create a new user
	user := User{ID: "1", Name: "Alice"}
	userJSON, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Check the response body
	var createdUser User
	json.Unmarshal(rr.Body.Bytes(), &createdUser)
	if createdUser.ID != user.ID || createdUser.Name != user.Name {
		t.Errorf("handler returned unexpected body: got %v want %v",
			createdUser, user)
	}
}

func TestGetUser(t *testing.T) {
	// Setup the handler and add a user to the map
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	users["1"] = User{ID: "1", Name: "Alice"}

	// Create a request to get the user
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var retrievedUser User
	json.Unmarshal(rr.Body.Bytes(), &retrievedUser)
	if retrievedUser.ID != "1" || retrievedUser.Name != "Alice" {
		t.Errorf("handler returned unexpected body: got %v want %v",
			retrievedUser, users["1"])
	}
}

func TestGetUserNotFound(t *testing.T) {
	// Setup the handler
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")

	// Create a request to get a non-existing user
	req, err := http.NewRequest("GET", "/users/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}
