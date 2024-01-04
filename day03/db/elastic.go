package db

import (
	"bytes"
	"context"
	"day03/types"
	"encoding/json"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

func PutData(id int, bulkIndexer esutil.BulkIndexer, data []byte) error {

	err := bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: strconv.Itoa(id),
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
	return err
}

func PutPlaces(places []types.Place, client *elasticsearch.Client) error {

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "place",
		Client: client,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, place := range places {
		data, err := json.Marshal(place)
		if err != nil {
			log.Fatal(err)
		}
		err = PutData(place.ID, bulkIndexer, data)
		if err != nil {
			return err
		}

	}

	err = bulkIndexer.Close(context.Background())

	return err
}
