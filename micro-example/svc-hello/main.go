package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	svcTimeEchoENVVAR string
	greetingPrefix    string
)

func livenessProbeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm alive!!!")
	slog.Info("Liveness probe performed")
}

func helloHandler(w http.ResponseWriter, _ *http.Request, svcTimeEchoENVVAR string, greetingPrefix string) {
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
	fmt.Fprintf(w, "%s at %s", greetingPrefix, string(respTime))

}

func Init() {
	var ok bool
	slog.Info("Load configuration")
	svcTimeEchoENVVAR, ok = os.LookupEnv("SVC_TIME_ECHO_URL")
	if !ok {
		slog.Error("Missing env variable SVC_TIME_ECHO_URL")
		os.Exit(1)
	}

	greetingPrefix, ok = os.LookupEnv("GREETING")
	if !ok {
		slog.Info("Missing env variable GREETING, used default value")
		greetingPrefix = "Hello world"
	}
}

func main() {
	var cancel context.CancelFunc

	Init()

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		helloHandler(w, r, svcTimeEchoENVVAR, greetingPrefix)
	})

	http.HandleFunc("/healthz", livenessProbeHandler)

	// Define the server
	srv := &http.Server{
		Addr: ":8890",
	}

	// Channel to listen for interrupt or terminate signals from the OS
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	slog.Info(fmt.Sprintf("Starting server on %s", srv.Addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Time Echo Service finished with an error. ", slog.String("error", err.Error()))
		}
	}()

	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				Init()
				slog.Info("SIGHUP recieved, reloading configuration...")

			// kill -SIGINT XXXX or Ctrl+c, kill -SIGTERM XXXX, kill -SIGQUIT XXXX
			case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				// Shutdown the server gracefully with a timeout context
				slog.Info("Shutting down server...")

				var ctx context.Context
				ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)

				if err := srv.Shutdown(ctx); err != nil {
					slog.Error("Server Shutdown Failed. ", slog.String("error", err.Error()))
				}
				cancel()
				exit_chan <- 0

			default:
				slog.Error("Unknown signal.")
			}
		}
	}()

	code := <-exit_chan
	os.Exit(code)
}
