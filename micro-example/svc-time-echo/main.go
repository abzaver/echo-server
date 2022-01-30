package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func livenessProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm alive!!!")
	slog.Info("Liveness probe performed")
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Received request to /time endpoint")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, time.Now().String())
}

func main() {
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/healthz", livenessProbeHandler)

	slog.Info("Starting Time Echo Service on port 8891")

	// Define the server
	srv := &http.Server{
		Addr: ":8891",
	}

	// Channel to listen for interrupt or terminate signals from the OS
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Goroutine to start the server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Time Echo Service finished with an error. ", slog.String("error", err.Error()))
		}
	}()

	// Readiness probe
	if _, err := os.Stat("/tmp"); os.IsNotExist(err) {
		if err = os.MkdirAll("/tmp", 0666); err != nil {
			slog.Error("Can't create temporary directory for health check file. ", slog.String("error", err.Error()))
		}
	}
	f, _ := os.Create("/tmp/healthy")
	_, err := f.WriteString("I'm alive!!!")
	if err != nil {
		slog.Error("Can't create temporary file for health check. ", slog.String("error", err.Error()))
	}
	f.Close()
	defer os.Remove(f.Name())

	slog.Info("Server is ready to handle requests at 0.0.0.0:8891")

	// Blocking main goroutine until an OS signal is received
	<-stop

	// Shutdown the server gracefully with a timeout context
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown Failed. ", slog.String("error", err.Error()))
	}

	slog.Info("Finished...")
}
