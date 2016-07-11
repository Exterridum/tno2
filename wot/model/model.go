package model

import "encoding/json"

type Context []interface{}

type Model struct {
	AT_Context Context    `json:"@context"`
	AT_Type    string     `json:"@type"`
	Name       string     `json:"name"`
	Uris       []string   `json:"uris"`
	Encodings  []string   `json:"encodings"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Name      string    `json:"name"`
	ValueType ValueType `json:"valueType"`
	Writable  bool      `json:"writable"`
	Hrefs     []string  `json:"hrefs"`
}

type ValueType struct {
	Type    string `json:"type"`
	Minimum int    `json:"minimum"`
	Maximum int    `json:"maximum"`
}

func (m Model) ToString() string {
	out, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return string(out)
}
