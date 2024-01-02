package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	var customers []Customer

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

	cert, err := os.ReadFile(*fCert)
	if err != nil {
		log.Fatal(err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{*fHost},
		Username:  *fUser,
		Password:  *fPassword,
		CACert:    cert,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "customers",
		Client: client,
	})
	if err != nil {
		log.Fatal(err)
	}

	names := []string{"Alice", "Ann", "Jhon", "Lisa", "Sophie", "Robert"}
	for i, n := range names {
		customers = append(customers, Customer{
			ID:   i,
			Name: n,
			Age:  rand.Intn(100),
		})
	}

	res, err := client.Indices.Create("customers")
	if err != nil {
		log.Fatalf("Cannot create index: %s\n", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s\n", res)
	}
	res.Body.Close()

	for _, customer := range customers {
		data, err := json.Marshal(customer)
		if err != nil {
			log.Fatal(err)
		}

		err = bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.Itoa(customer.ID),
			Body:       bytes.NewReader(data),
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				log.Println("Successful added item")
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("ERROR: %s", err)
				} else {
					log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
				}
			},
		})

		if err != nil {
			log.Fatal(err)
		}

	}

	if err := bulkIndexer.Close(context.Background()); err != nil {
		log.Fatal(err)
	}

}
