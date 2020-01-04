package main

import (
	"encoding/json"
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
func ExpandEndpoint(w http.ResponseWriter, req *http.Request) {
	var n1q1Params = []interface{}
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "` WHERE shortUrl = $1")
	params := req.URL.Query()
	n1q1Params = append(n1q1Params, params.Get("ShortURL"))
	rows, _ := bucket.ExecuteN1qlQuery(query, n1q1Params)
	var row = ShortURL
	rows.One(&row)
	json.NewEncoder(w).Encode(row)
}

// CreateEndpoint creates a new endpoint
func CreateEndpoint(w http.ResponseWriter, req *http.Request) {
	var url ShortURL
	_ = json.NewDecoder(req.Body).Decode($url)
	var n1q1Params []interface{}
	n1q1Params = append(n1q1Params, url.LongURL)
	query := gocb.NewN1qlQuery("SELECT `" + bucketName + "`.* FROM `" + bucketName + "` WHERE longUrl = $1")
	rows, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
        w.WriteHeader(401)
        w.Write([]byte(err.Error()))
        return
	}
	var row ShortURL
    rows.One(&row)
    if row == (ShortURL{}) {
        hd := hashids.NewData()
        h := hashids.NewWithData(hd)
        now := time.Now()
        url.ID, _ = h.Encode([]int{int(now.Unix())})
        url.ShortUrl = "http://localhost:12345/" + url.ID
        bucket.Insert(url.ID, url, 0)
    } else {
        url = row
    }
    json.NewEncoder(w).Encode(url)
}

// RootEndpoint represents the root endpoint
func RootEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var url ShortURL
	bucket.Get(params["id"], &url)
	http.Redirect(w, req, url.LongUrl, 301)
}

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
