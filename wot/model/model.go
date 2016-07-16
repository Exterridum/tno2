package model

import "encoding/json"

type Context []interface{}

type ThingDescription struct {
	AT_Context Context    `json:"@context"`
	AT_Type    string     `json:"@type"`
	Name       string     `json:"name"`
	Uris       []string   `json:"uris"`
	Encodings  []string   `json:"encodings"`
	Properties []Property `json:"properties"`
	Actions    []Action   `json:"actions"`
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

func (td *ThingDescription) ToString() string {
	out, err := json.Marshal(td)
	if err != nil {
		panic(err)
	}

	return string(out)
}
