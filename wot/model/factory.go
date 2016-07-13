package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//TODO: Implement extension pattern
type Loader interface {
	Load(url string) ThingDescription
}

func Load(uri string) *ThingDescription {
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
