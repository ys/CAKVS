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

type CakvsENV struct{}
type CakvsVal map[string]string

func (m *CakvsENV) Get(k string) CakvsVal {
	var dat CakvsVal
	val := os.Getenv(k)
	json.Unmarshal([]byte(val), &dat)
	return dat
}

func (m *CakvsENV) Set(k string, val *CakvsVal) {
	valString, _ := json.Marshal(val)
	os.Setenv(k, string(valString))
}

func (m *CakvsENV) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		var val CakvsVal
		if err := json.Unmarshal(body, &val); err != nil {
			panic(err)
		}
		fmt.Println(string(body), key, val)
		m.Set(key, &val)
	}
}

func main() {
	m := CakvsENV{}
	http.Handle("/", &m)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
