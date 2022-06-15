package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	clientv3 "go.etcd.io/etcd/client/v3"
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

func registerServiceWithEtcd() {
	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to connect to etcd: %v", err)
	}
	defer cli.Close()

	serviceKey := "/services/user-service"
	serviceValue := "localhost:8000"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = cli.Put(ctx, serviceKey, serviceValue)
	cancel()
	if err != nil {
		log.Fatalf("Failed to register service with etcd: %v", err)
	}

	log.Printf("Service registered with etcd: %s -> %s", serviceKey, serviceValue)
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

	go registerServiceWithEtcd()

	log.Println("User Service is running on port 8000")
	http.ListenAndServe(":8000", router)
}
