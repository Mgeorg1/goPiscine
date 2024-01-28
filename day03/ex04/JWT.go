package main

import (
	"day03/db"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secretKey")

func createToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString(jwtKey)
}

type JsonTokenResponse struct {
	Token string `json:"token"`
}

func GetJWTTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	token, err := createToken()
	if err != nil {
		log.Printf("generating JWT token error: %s", err)
		db.WriteErrorJson(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonToken := JsonTokenResponse{Token: token}
	db.WriteJSON(w, jsonToken)
}

func verifyJWTToken(endpointHandler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authValue := r.Header.Get("Authorization")
		if authValue == "" {
			db.WriteErrorJson(w, "Access token not found", http.StatusUnauthorized)
			log.Printf("client unauthorized, client ip: %s", r.Host)
			return
		}

		tokenSlice := strings.Split(authValue, " ")
		if len(tokenSlice) != 2 || tokenSlice[0] != "Bearer" {
			db.WriteErrorJson(w, "Invalid Authorization header format", http.StatusUnauthorized)
			log.Printf("client unauthorized, client ip: %s", r.Host)
			return
		}

		token, err := jwt.Parse(tokenSlice[1], func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			db.WriteErrorJson(w, "invalid token", http.StatusUnauthorized)
			log.Printf("client unauthorized, client ip: %s", r.Host)
			return
		}
		endpointHandler(w, r)
	}
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
	server := db.CreateServer(store)
	http.HandleFunc("/api/get_token", GetJWTTokenHandler)
	http.HandleFunc("/api/recommend", verifyJWTToken(server.GetClosestPlacesHandler))
	err = http.ListenAndServe(":8800", nil)
	if err != nil {
		log.Fatal(err)
	}

}
