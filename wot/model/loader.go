package device

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//TODO: Implement extension pattern
type Loader interface {
	Load(url string) Model
}

func Load(uri string) Model {
	sep := strings.Split(uri, "://")
	method, path := sep[0], sep[1]

	if method == "file" {
		return fromFile(path)
	}

	return Model{}
}

func fromFile(path string) Model {
	file, e := ioutil.ReadFile(path)

	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var model Model

	json.Unmarshal(file, &model)

	return model
}
