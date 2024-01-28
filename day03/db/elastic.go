package db

import (
	"bytes"
	"context"
	"day03/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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

func (elasticStore *ElasticStore) СreateIndex(indexName string) {
	es := elasticStore.Client
	mapping := `
	{
		"mappings": {
			"properties": {
				"name": {
				  "type":  "text"
				},
				"address": {
				  "type":  "text"
				},
				"phone": {
				  "type":  "text"
				},
				"location": {
				  "type": "geo_point"
				}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error creating index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error creating index: %s", res.String())
	} else {
		log.Printf("Index %s created successfully", indexName)
	}
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

	var places []types.Place
	for _, hit := range responses.Hits.Hits {
		places = append(places, hit.Hit)
	}
	return places, responses.Hits.Total.Value, nil
}

type Server struct {
	store *ElasticStore
}

func CreateServer(store *ElasticStore) *Server {
	return &Server{store: store}
}

type ResponsePlaces struct {
	Places       []types.Place
	Total        int
	IsPrevious   bool
	IsNext       bool
	PreviousPage int
	NextPage     int
	LastPage     int
}

type JsonClosestResponse struct {
	Name   string        `json:"name"`
	Places []types.Place `json:"places"`
}

type JsonResponsePlaces struct {
	Name     string        `json:"name"`
	Total    int           `json:"total"`
	Places   []types.Place `json:"places"`
	PrevPage int           `json:"prev_page"`
	NextPage int           `json:"next_page"`
	LastPage int           `json:"last_page"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetErrorResp(message string) []byte {
	resp := ErrorResponse{Error: message}
	errorJSON, _ := json.Marshal(resp)
	return errorJSON
}

func (server *Server) getPagePlaces(w http.ResponseWriter, r *http.Request) (ResponsePlaces, error) {
	var resp ResponsePlaces
	pageParam := r.URL.Query().Get("page")
	var page int
	var err error
	if pageParam == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageParam)
	}
	if err != nil {
		log.Printf("Error while converting string 'page' value. Error: %s\n", err)
		WriteErrorJson(w, "Error while converting string 'page' value "+err.Error(), http.StatusInternalServerError)
		return resp, err
	}
	limit := 10
	offset := (page - 1) * 10
	if offset < 0 {
		log.Printf("Error while converting string 'page' value. Error: %s\n", err)

		WriteErrorJson(w, "offset < 0", http.StatusBadRequest)
		err = fmt.Errorf("offset for places page < 0")
		return resp, err
	}
	places, hitNum, err := server.store.GetPlaces(limit, page)
	if err != nil {
		log.Println(err)
		WriteErrorJson(w, "Places not found", http.StatusNotFound)
		return resp, err
	}

	resp.Places = places
	resp.Total = hitNum

	if offset > 0 {
		resp.IsPrevious = true
		resp.PreviousPage = page - 1
	}

	if offset+limit < hitNum {
		resp.IsNext = true
		resp.NextPage = page + 1
	}

	resp.LastPage = (hitNum + limit - 1) / limit

	if page > resp.LastPage {
		WriteErrorJson(w, "Invalid 'page' value "+pageParam, http.StatusBadRequest)
		return resp, err
	}
	return resp, nil
}

func (server *Server) GetPlacesHandler(w http.ResponseWriter, r *http.Request) {

	data, err := server.getPagePlaces(w, r)
	if err != nil {
		return
	}
	tmpl := template.Must(template.ParseFiles("template.html"))

	err = tmpl.Execute(w, data)
	if err != nil {
		WriteErrorJson(w, "Error executing template"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func WriteErrorJson(w http.ResponseWriter, errorText string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(GetErrorResp(errorText))
}

func WriteJSON(w http.ResponseWriter, response interface{}) {
	json, err := json.MarshalIndent(response, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		WriteErrorJson(w, "Error due JSON marhalling: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(json)
	if err != nil {
		WriteErrorJson(w, "Error while writing response header: "+err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) GetClosestPlaces(lat float64, lon float64) ([]types.Place, error) {
	query := map[string]interface{}{
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"order":           "asc",
					"unit":            "km",
					"mode":            "min",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	queryString := string(queryBytes)
	log.Println(queryString)

	client := server.store.Client
	res, err := client.Search(client.Search.WithIndex("place"),
		client.Search.WithSize(3),
		client.Search.WithBody(strings.NewReader(queryString)))
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		return nil, errors.New(res.Status())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responses types.Responses
	var places []types.Place
	err = json.Unmarshal(body, &responses)

	if err != nil {
		return nil, err
	}
	for _, hit := range responses.Hits.Hits {
		places = append(places, hit.Hit)
	}

	return places, nil
}

func (server *Server) GetClosestPlacesHandler(w http.ResponseWriter, r *http.Request) {
	latParam := r.URL.Query().Get("lat")
	lonParam := r.URL.Query().Get("lon")

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		WriteErrorJson(w, err.Error()+"lat param: "+latParam, http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonParam, 64)
	if err != nil {
		WriteErrorJson(w, err.Error()+"lon param: "+lonParam, http.StatusBadRequest)
		return
	}

	places, err := server.GetClosestPlaces(lat, lon)

	if err != nil {
		log.Printf("error when getting closes places %s", err)
		WriteErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := JsonClosestResponse{
		Name:   "Recomendation",
		Places: places,
	}
	WriteJSON(w, response)
}

func (server *Server) GetJsonPlacesHandler(w http.ResponseWriter, r *http.Request) {
	var jsonResp interface{}
	data, err := server.getPagePlaces(w, r)
	if err != nil {
		return
	}
	jsonResp = JsonResponsePlaces{
		Name:     "Places",
		Total:    data.Total,
		Places:   data.Places,
		PrevPage: data.PreviousPage,
		NextPage: data.NextPage,
		LastPage: data.LastPage,
	}

	WriteJSON(w, jsonResp)
}
