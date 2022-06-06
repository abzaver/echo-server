package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = make(map[string]User)

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	users[user.ID] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func registerServiceWithConsul() {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	serviceID := "user-service"
	address := "localhost"
	port := 8000

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "user-service",
		Address: address,
		Port:    port,
		Tags:    []string{"primary"},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
			Interval: "10s",
		},
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Failed to register service with Consul: %v", err)
	}

	log.Printf("Service registered with Consul: %s", serviceID)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")

	go registerServiceWithConsul()

	log.Println("User Service is running on port 8000")
	http.ListenAndServe(":8000", router)
}
