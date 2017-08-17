package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	tribHandler := NewTribHandler()

	// routing
	router := mux.NewRouter()
	router.HandleFunc("/tribonacci/{n}", tribHandler.TribonacciHandler).Methods("GET")

	// http server
	h := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		log.Info("Start listening on :8080")

		if err := h.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-stop

	log.Info("Shutting down the server...")
	ctx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	err := h.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
