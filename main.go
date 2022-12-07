package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tm := time.Now().Format(time.RFC1123)
		_, _ = w.Write([]byte("The time is: " + tm))
	})

	// Set up the HTTP server:
	server := &http.Server{}
	server.Addr = ":8080"
	server.Handler = mux
	server.SetKeepAlivesEnabled(true)
	server.ReadTimeout = 15 * time.Minute
	server.WriteTimeout = 15 * time.Minute

	// Start the server:
	logrus.Printf("listening on address: %v", server.Addr)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			//Normal graceful shutdown error
			if err.Error() == "http: Server closed" {
				logrus.Info(err)
			} else {
				logrus.Fatal(err)
			}
		}
	}()
	//Wait for shutdown signal, then shutdown api server. This will wait for all connections to finish.
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interruptChan
	logrus.Debug("Shutting down API server...")

	err := server.Shutdown(context.Background())
	if err != nil {
		logrus.Error("Error shutting down server: ", err)
	}
	logrus.Info("terminated")
}
