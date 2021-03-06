package diagnostics

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewDiagnostics() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", health)
	router.HandleFunc("/ready", ready)

	return router
}

func health(w http.ResponseWriter, r *http.Request) {
	log.Print("The health handler was called.")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}

func ready(w http.ResponseWriter, r *http.Request) {
	log.Print("The ready handler was called.")
	fmt.Fprint(w, http.StatusText(http.StatusOK))
}
