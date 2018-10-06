package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kuskmen/golang-playground/internal/diagnostics"
)

type ServerConfig struct {
	port   string
	router http.Handler
	name   string
}

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

	select {
	case err := <-errorChannel:
		context, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		for _, s := range servers {
			shutDownError := s.Shutdown(context)
			if shutDownError != nil {
				log.Print(shutDownError)
			}

			log.Println("Server gracefully shutdown.")
		}

		log.Fatal(err)
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Print("The hello handler was called.")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
