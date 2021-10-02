package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to /time endpoint")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, time.Now().String())
	})

	log.Println("Starting Time Echo Service on port 8891")

	err := http.ListenAndServe(":8891", nil)
	if err != nil {
		log.Fatalln("Time Echo Service finished with an error.", err)
	}
}
