package model

import "encoding/json"

type Context []interface{}

type Model struct {
	AT_Context Context  `json:"@context"`
	AT_Type    string   `json:"@type"`
	Name       string   `json:"name"`
	Uris       []string `json:"uris"`
	Encodings  []string `json:"encodings"`
}

type Property struct {
	Name      string    `json: "Name"`
	ValueType ValueType `json: "valueType"`
	Writable  bool      `json: "writable"`
	Hrefs     []string  `json: "hrefs"`
}

type ValueType struct {
	Type string `json: "type"`
}

func (m Model) ToString() string {
	out, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(out)
}
