package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	svcTimeEchoENVVAR, ok := os.LookupEnv("SVC_TIME_ECHO_URL")
	if !ok {
		log.Fatalln("Missing env variable SVC_TIME_ECHO_URL")
		os.Exit(1)
	}

	svcTimeEchoURL := fmt.Sprintf("%s/time", svcTimeEchoENVVAR)
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request to /hello endpoint")
		svcTimeEchoResp, err := http.DefaultClient.Get(svcTimeEchoURL)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Service TimeEcho (%s) is not available: %v", svcTimeEchoURL, err)
			return
		}
		defer svcTimeEchoResp.Body.Close()
		respTime, _ := io.ReadAll(svcTimeEchoResp.Body)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello World at %s", string(respTime))
	})
	log.Println("Starting server on port 8890")
	err := http.ListenAndServe(":8890", nil)
	if err != nil {
		log.Fatalln("Hello Server finished with an error.", err.Error())
	}
}
