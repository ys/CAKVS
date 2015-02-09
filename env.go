package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	m := MongoENV{}
	m.Set("LOL", &MongoVal{"YOLO": "TRUEFRIEND"})
	fmt.Println(m.Get("LOL"))
	fmt.Println(m.Get("DAT"))
}
