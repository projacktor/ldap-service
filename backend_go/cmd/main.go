package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	log.Println("Server is running")

	if err := http.ListenAndServe(":3001", router); err != nil {
		log.Fatal(err)
	}

}
