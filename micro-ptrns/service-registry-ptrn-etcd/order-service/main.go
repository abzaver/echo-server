package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

func getUserServiceAddress(cli *clientv3.Client) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := cli.Get(ctx, "/services/user-service")
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("user service not found")
	}

	return string(resp.Kvs[0].Value), nil
}

func getUserServiceAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Connect to etcd for User Service search
	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		http.Error(w, "Failed to connect to etcd", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	userServiceAddr, err := getUserServiceAddress(cli)
	if err != nil {
		http.Error(w, "User Service not found", http.StatusInternalServerError)
		return
	}

	userInfoResp, err := http.Get(fmt.Sprintf("http://%s/health", userServiceAddr))
	if err != nil || userInfoResp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to get health info. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(userServiceAddr))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	json.NewDecoder(r.Body).Decode(&order)
	orders[order.ID] = order
	// Подключение к etcd для поиска User Service
	etcdEndpoint := os.Getenv("ETCD_ENDPOINT")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		http.Error(w, "Failed to connect to etcd", http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	userServiceAddr, err := getUserServiceAddress(cli)
	if err != nil {
		http.Error(w, "User Service not found", http.StatusInternalServerError)
		return
	}

	// Запрос к User Service для получения информации о пользователе
	log.Printf("http://%s/users/%s", userServiceAddr, order.UserID)
	userInfoResp, err := http.Get(fmt.Sprintf("http://%s/users/%s", userServiceAddr, order.UserID))
	if err != nil || userInfoResp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Failed to get user info. Status code: %d", userInfoResp.StatusCode), http.StatusInternalServerError)
		return
	}

	var userInfo map[string]interface{}
	json.NewDecoder(userInfoResp.Body).Decode(&userInfo)

	// Теперь у нас есть информация о пользователе, и можно продолжить обработку заказа
	orderInfo := map[string]interface{}{
		"order": order,
		"user":  userInfo,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderInfo)
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
	serviceValue := "order-service:8001"

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
	err := godotenv.Load("../.env")
	if err != nil {
		log.Printf("No .env file found, relying on system environment variables")
	}

	router := mux.NewRouter()
	router.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	router.HandleFunc("/orders", createOrder).Methods("POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")
	router.HandleFunc("/getuserserviceaddress", getUserServiceAddressHandler).Methods("GET")

	go registerServiceWithEtcd()

	log.Println("Order Service is running on port 8001")
	http.ListenAndServe(":8001", router)
}
