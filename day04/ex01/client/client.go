package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"day04/api"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	candyType := flag.String("k", "", "candy type")
	candyCount := flag.Int("c", 0, "candy count")
	money := flag.Int("m", 0, "money amount")
	flag.Parse()

	if candyType == nil || *candyType == "" {
		log.Fatalln("candy type is required")
	}
	req := api.BuyCandyRequest{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *candyCount,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	cert, err := tls.LoadX509KeyPair("../client.localhost/cert.pem", "../client.localhost/key.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCert, err := os.ReadFile("../minica.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
	r := bytes.NewReader(reqJson)
	resp, err := client.Post("https://localhost:3333/buy_candy", "application/json", r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var buyCandyResp api.BuyCandyResponse
	err = json.NewDecoder(resp.Body).Decode(&buyCandyResp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s Your change: %d", buyCandyResp.Thanks, buyCandyResp.Change)
}
