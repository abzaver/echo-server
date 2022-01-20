package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

func livenessProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm alive!!!")
	slog.Info("Liveness probe performed")
}

func helloHandler(w http.ResponseWriter, _ *http.Request, svcTimeEchoENVVAR string) {
	slog.Info("Received request to /hello endpoint")
	svcTimeEchoURL := fmt.Sprintf("%s/time", svcTimeEchoENVVAR)
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

}

func main() {
	svcTimeEchoENVVAR, ok := os.LookupEnv("SVC_TIME_ECHO_URL")
	if !ok {
		slog.Error("Missing env variable SVC_TIME_ECHO_URL")
		os.Exit(1)
	}

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		helloHandler(w, r, svcTimeEchoENVVAR)
	})

	http.HandleFunc("/healthz", livenessProbeHandler)

	slog.Info("Starting server on port 8890")
	err := http.ListenAndServe(":8890", nil)
	if err != nil {
		slog.Error("Hello Server finished with an error.", slog.String("error", err.Error()))
	}
}
