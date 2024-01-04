package main

import (
	"day03/db"
	"day03/types"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

func parseCsv(filePath string) ([]types.Place, error) {
	var ret []types.Place

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
		var place types.Place

		place.ID, err = strconv.Atoi(row[0])
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[0], err)
			continue
		}
		place.Name = row[1]
		place.Address = row[2]
		place.Phone = row[3]
		place.Loc.Lon, err = strconv.ParseFloat(row[4], 64)
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[4], err)
			continue
		}
		place.Loc.Lat, err = strconv.ParseFloat(row[5], 64)
		if err != nil {
			log.Printf("Converting error. Skip row. Column: %s, error: %s", row[5], err)
			continue
		}

		ret = append(ret, place)
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

	res, err := client.Indices.Create("place")
	if err != nil {
		log.Fatalf("Cannot create index: %s\n", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s\n", res)
	}
	res.Body.Close()

	places, err := parseCsv(*fData)
	if err != nil {
		log.Fatal(err)
	}

	err = db.PutPlaces(places, client)
	if err != nil {
		log.Fatal(err)
	}
}
