package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func routesHandler(w http.ResponseWriter, r *http.Request) {
	source, ok := r.URL.Query()["src"]
	if !ok || len(source) != 1 {
		handleError(w, http.StatusBadGateway, "one src parameter must be specified")
		return
	}

	_, ok = r.URL.Query()["dst"]
	if !ok {
		handleError(w, http.StatusBadGateway, "at least one dst parameter must be specified")
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("ok")
}

func handleError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(message)
}

func main() {
	// TODO handle not found
	http.HandleFunc("/routes", routesHandler)

	log.Println("starting listener")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// TODO handle graceful shutdown
}
