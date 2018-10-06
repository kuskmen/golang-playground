package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/kuskmen/golang-playground/internal/diagnostics"
)

type ServerConfig struct {
	port   string
	router http.Handler
	name   string
}

var requestsCount uint64

func main() {
	log.Print("Starting the application...")

	/* 	blPort := os.Getenv("PORT")
	   	if len(blPort) == 0 {
	   		log.Fatal("The application port should be set")
	   	}

	   	diagPort := os.Getenv("DIAG_PORT")
	   	if len(diagPort) == 0 {
	   		log.Fatal("The diagnostics port should be set")
	   	} */

	router := mux.NewRouter()
	router.HandleFunc("/", hello)
	diagRouter := diagnostics.NewDiagnostics()

	errorChannel := make(chan error, 2)

	configs := []ServerConfig{
		{
			port:   "8080",
			router: router,
			name:   "Application server",
		},
		{
			port:   "8585",
			router: diagRouter,
			name:   "Diagnostics server",
		},
	}

	servers := make([]*http.Server, 2)

	for _, c := range configs {
		i := 0
		go func(config ServerConfig, i int) {
			log.Printf("The %s is preparing to handle connections...", config.name)

			servers[i] = &http.Server{
				Addr:    ":" + config.port,
				Handler: config.router,
			}

			err := servers[i].ListenAndServe()
			if err != nil {
				errorChannel <- err
			}
		}(c, i)
		i++
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errorChannel:
		log.Printf("Received an error: %v", err)
	case sig := <-interrupt:
		log.Printf("Received the signal %v", sig)
	}

	for _, s := range servers {
		context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := s.Shutdown(context)
		if err != nil {
			log.Print(err)
		}
		log.Println("Server gracefully shutdown.")
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called.")

	atomic.AddUint64(&requestsCount, 1)
	fmt.Fprintf(w, "The hello handler was called %v times.", atomic.LoadUint64(&requestsCount))
}
