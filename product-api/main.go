package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/product-api/handlers"
)

func main() {
	l := log.New(os.Stdout, "product_api_", log.LstdFlags)
	sm := http.NewServeMux()
	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	s := &http.Server{
		Addr:         "localhost:9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, os.Kill)
	sig := <-sigchan
	l.Println("Received terminate, Graceful shutdown", sig)
	ct, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ct)
}
