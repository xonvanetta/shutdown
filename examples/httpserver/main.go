package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/xonvanetta/shutdown/pkg/shutdown"
)

func main() {
	server := http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to listen to server: %s", err)
		}
	}()

	<-shutdown.Chan()
	log.Println("shutdown was called")
	ctx, cancel := context.WithCancel(context.TODO())
	go func() {
		time.Sleep(time.Second * 30)
		cancel()
	}()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("failed to do graceful shutdown for given time: %s", err)
	}
	log.Println("shutdown was finished")
}
