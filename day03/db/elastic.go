package db

import (
	"bytes"
	"context"
	"day03/types"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
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

type ElasticStore struct {
	Client *elasticsearch.Client
}

func (elasticStore *ElasticStore) PutPlaces(places []types.Place) error {

	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "place",
		Client: elasticStore.Client,
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

func CreateElasticStore(certFilePath, host, user, passw *string) (*ElasticStore, error) {
	cert, err := os.ReadFile(*certFilePath)
	if err != nil {
		return nil, err
	}

	cfg := elasticsearch.Config{
		Addresses: []string{*host},
		Username:  *user,
		Password:  *passw,
		CACert:    cert,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	es := new(ElasticStore)
	es.Client = client
	return es, nil
}

func (elasticStore *ElasticStore) GetPlaces(limit int, offset int) ([]types.Place, int, error) {
	resp, err := elasticStore.Client.Search(elasticStore.Client.Search.WithIndex("place"),
		elasticStore.Client.Search.WithFrom(offset),
		elasticStore.Client.Search.WithSize(limit))
	if err != nil {
		return nil, 0, err
	}

	if resp.IsError() {
		return nil, 0, errors.New(resp.Status())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var responses types.Responses

	err = json.Unmarshal(body, &responses)
	if err != nil {
		return nil, 0, err
	}
	log.Println(string(body))
	var places []types.Place
	for _, hit := range responses.Hits.Hits {
		places = append(places, hit.Hit)
	}
	return places, responses.Hits.Total.Value, nil
}
