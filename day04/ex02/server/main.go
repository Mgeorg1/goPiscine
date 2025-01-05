package main

import (
	"C"
	"crypto/tls"
	"crypto/x509"
	"day04/api"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buy_candy", api.BuyCandyHandlerCow)

	caCert, err := os.ReadFile("../../ex01/minica.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	server := &http.Server{
		Addr:      "localhost:3333",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	err = server.ListenAndServeTLS("../../ex01/localhost/cert.pem", "../../ex01/localhost/key.pem")
	if err != nil {
		log.Fatalln(err)
	}
}
