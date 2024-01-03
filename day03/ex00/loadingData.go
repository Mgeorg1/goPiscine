package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type Restaurant struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Phone   string   `json:"phone"`
	Loc     Location `json:"location"`
}

type Location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func parseCsv(filePath string) ([]Restaurant, error) {
	var ret []Restaurant

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, row := range data {
		var rest Restaurant

		rest.ID, err = strconv.Atoi(row[0])
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[0], err)
			continue
		}
		rest.Name = row[1]
		rest.Address = row[2]
		rest.Phone = row[3]
		rest.Loc.Lon, err = strconv.ParseFloat(row[4], 64)
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[4], err)
			continue
		}
		rest.Loc.Lat, err = strconv.ParseFloat(row[5], 64)
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[5], err)
			continue
		}

		ret = append(ret, rest)
	}
	log.Printf("Generated data len: %d\n", len(ret))
	return ret, nil
}

func main() {

	log.SetFlags(0)
	fHost := flag.String("h", "https://localhost:9200", " -h host:port")
	fCert := flag.String("cacert", "./http_ca.crt", "-cacert ./path to ca certificate")
	fUser := flag.String("u", "", "-u username")
	fPassword := flag.String("p", "", "-p password")
	fData := flag.String("f", "", "-f csv file")

	flag.Parse()

	if *fUser == "" || *fPassword == "" || *fData == "" {
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
		Index:  "restaurants",
		Client: client,
	})
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Indices.Create("restaurants")
	if err != nil {
		log.Fatalf("Cannot create index: %s\n", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s\n", res)
	}
	res.Body.Close()

	restaurants, err := parseCsv(*fData)
	if err != nil {
		log.Fatal(err)
	}

	for _, rest := range restaurants {
		data, err := json.Marshal(rest)
		if err != nil {
			log.Fatal(err)
		}

		err = bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action:     "index",
			DocumentID: strconv.Itoa(rest.ID),
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
