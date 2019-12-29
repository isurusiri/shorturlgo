package main

import (
	"log"
	"net/http"

	"github.com/couchbase/gocb"
	"github.com/gorilla/mux"
)

var bucket *gocb.Bucket
var bucketName string

// ShortURL represents the properties of a short url
type ShortURL struct {
	ID       string `json:"id,omitempty"`
	LongURL  string `json:"longUrl,omitempty"`
	ShortURL string `json:"shortUrl,omitempty"`
}

// ExpandEndpoint exapands a new endpoint
func ExpandEndpoint(w http.ResponseWriter, req *http.Request) {}

// CreateEndpoint creates a new endpoint
func CreateEndpoint(w http.ResponseWriter, req *http.Request) {}

// RootEndpoint represents the root endpoint
func RootEndpoint(w http.ResponseWriter, req *http.Request) {}

func main() {
	router := mux.NewRouter()
	cluster, _ := gocb.Connect("couchbase://localhost")
	bucketName = "ISIRLK-1-TEST-1"
	bucket, _ = cluster.OpenBucket(bucketName, "")

	router.HandleFunc("/{id}", RootEndpoint).Methods("GET")
	router.HandleFunc("/expand/", ExpandEndpoint).Methods("GET")
	router.HandleFunc("/create", CreateEndpoint).Methods("PUT")

	log.Fatal(http.ListenAndServe(":12345", router))
}
