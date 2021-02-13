package main

import (
	"log"
	"net/http"

	"github.com/rafarlopes/route-service/internal/route"
)

func main() {
	// TODO handle not found
	http.HandleFunc("/routes", route.RoutesHandler)

	log.Println("starting listener")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// TODO handle graceful shutdown
}
