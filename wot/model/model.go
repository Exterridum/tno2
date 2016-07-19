package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Context []interface{}

type ThingDescription struct {
	AT_Context Context    `json:"@context"`
	AT_Type    string     `json:"@type"`
	Name       string     `json:"name"`
	Uris       []string   `json:"uris"`
	Encodings  []string   `json:"encodings"`
	Properties []Property `json:"properties"`
	Actions    []Action   `json:"actions"`
	Events     []Event    `json:"events"`
}

type Property struct {
	Name      string    `json:"name"`
	ValueType ValueType `json:"valueType"`
	Unit      string    `json:"unit"`
	Writable  bool      `json:"writable"`
	Hrefs     []string  `json:"hrefs"`
}

type Action struct {
	AT_Type    string     `json:"@type"`
	Name       string     `json:"name"`
	InputData  InputData  `json:"inputData"`
	OutputData OutputData `json:"outputData"`
	Hrefs      []string   `json:"hrefs"`
}

type Event struct {
	AT_Type   string    `json:"@type"`
	Name      string    `json:"name"`
	ValueType ValueType `json:"valueType"`
	Hrefs     []string  `json:"hrefs"`
}

type InputData struct {
	ValueType ValueType `json:"valueType"`
	Unit      string    `json:"unit"`
}

type OutputData struct {
	ValueType ValueType `json:"valueType"`
	Unit      string    `json:"unit"`
}

type ValueType struct {
	Type    string `json:"type"`
	Minimum int    `json:"minimum"`
	Maximum int    `json:"maximum"`
}

func Create(uri string) *ThingDescription {
	sep := strings.Split(uri, "://")
	method, path := sep[0], sep[1]

	if method == "file" {
		return fromFile(path)
	}

	return &ThingDescription{}
}

func fromFile(path string) *ThingDescription {
	file, e := ioutil.ReadFile(path)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var td ThingDescription

	json.Unmarshal(file, &td)

	td.Uris = make([]string, 0)

	return &td
}
