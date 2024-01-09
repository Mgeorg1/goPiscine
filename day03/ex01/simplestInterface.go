package main

import (
	"day03/db"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
	store *db.ElasticStore
}

func createServer(store *db.ElasticStore) *Server {
	return &Server{store: store}
}

func (server *Server) GetPlacesHandler(w http.ResponseWriter, r *http.Request) {
	// data, err := server.PrepareData(r)
	// if err.Message != "" {
	// 	http.Error(w, err.Message, http.StatusInternalServerError)
	// 	return
	// }

	data, count, _ := server.store.GetPlaces(2000, 1)
	var str string
	if len(data) != 0 {
		log.Println("!")
		str = fmt.Sprintf("%s\n%s\n%s\n", data[0].Name, data[0].Address, data[0].Phone)
	}
	log.Println(str)
	fmt.Fprintf(w, "Sorry, it's not finished, but this is hit count num: %d\nFirst place:\n%s", count, str)
	// tmpl := template.Must(template.ParseFiles("web/template/template.html"))
	// if err := tmpl.Execute(w, data); err != nil {
	// 	http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	// }
}

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello! Type place to index")
}

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
	server := createServer(store)
	http.HandleFunc("/place", server.GetPlacesHandler)
	http.HandleFunc("/", Hello)
	err = http.ListenAndServe(":8800", nil)
	if err != nil {
		log.Fatal(err)
	}

}
