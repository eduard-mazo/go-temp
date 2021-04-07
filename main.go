package main

import (
	"log"
	"net/http"
	"os"

	"example.com/hello/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.Grettings).Methods("GET")
	r.HandleFunc("/sensor/{sensorID}/temp/{temp}", handlers.UpdateTemp).Methods("GET")
	r.HandleFunc("/sensor/{sensorID}", handlers.GetTemp).Methods("GET")
	// r.HandleFunc("/", deleteSensor).Methods("DELETE")

	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
