package main

import (
	"day03/db"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(0)
	fHost := flag.String("h", "https://localhost:9200", " -h host:port")
	fCert := flag.String("cacert", "./http_ca.crt", "-cacert ./path to ca certificate")
	fUser := flag.String("u", "", "-u username")
	fPassword := flag.String("p", "", "-p password")

	flag.Parse()

	if *fUser == "" || *fPassword == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	store, err := db.CreateElasticStore(fCert, fHost, fUser, fPassword)
	if err != nil {
		log.Fatal(err)
	}
	server := db.CreateServer(store)
	http.HandleFunc("/api/places", server.GetJsonPlacesHandler)
	err = http.ListenAndServe(":8800", nil)
	if err != nil {
		log.Fatal(err)
	}

}
