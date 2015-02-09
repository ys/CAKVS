package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type MongoENV struct{}
type MongoVal map[string]string

func (m *MongoENV) Get(k string) MongoVal {
	var dat MongoVal
	val := os.Getenv(k)
	json.Unmarshal([]byte(val), &dat)
	return dat
}

func (m *MongoENV) Set(k string, val *MongoVal) {
	valString, _ := json.Marshal(val)
	os.Setenv(k, string(valString))
}

func (m *MongoENV) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/"):]
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if r.Method == "GET" {
		if err := json.NewEncoder(w).Encode(m.Get(key)); err != nil {
			panic(err)
		}
	}
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			panic(err)
		}
		if err := r.Body.Close(); err != nil {
			panic(err)
		}
		var val MongoVal
		if err := json.Unmarshal(body, &val); err != nil {
			panic(err)
		}
		fmt.Println(string(body), key, val)
		m.Set(key, &val)
	}
}

func main() {
	m := MongoENV{}
	http.Handle("/", &m)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
