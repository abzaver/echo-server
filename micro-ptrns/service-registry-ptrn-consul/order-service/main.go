package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hashicorp/consul/api"
)

type Order struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
}

var orders = make(map[string]Order)

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	order, exists := orders[id]
	if !exists {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(order)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	orders[order.ID] = order
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func registerServiceWithConsul() {
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	serviceID := "order-service"
	address := "localhost"
	port := 8001

	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    "order-service",
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
	router.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")

	go registerServiceWithConsul()

	log.Println("Order Service is running on port 8001")
	http.ListenAndServe(":8001", router)
}
