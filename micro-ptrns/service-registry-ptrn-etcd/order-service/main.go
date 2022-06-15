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

	serviceKey := "/services/order-service"
	serviceValue := "localhost:8001"

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
	router.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")

	go registerServiceWithEtcd()

	log.Println("Order Service is running on port 8001")
	http.ListenAndServe(":8001", router)
}
