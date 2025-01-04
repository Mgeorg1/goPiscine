package main

import (
	"day04/api"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buy_candy", api.BuyCandyHandler)
	err := http.ListenAndServe("0.0.0.0:3333", mux)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
