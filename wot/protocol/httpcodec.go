package protocol

import (
	"encoding/json"
	"net/http"

	"github.com/conas/tno2/wot/model"
)

type WotNumber struct {
	Value float64 `json:"value"`
}

type WotString struct {
	Value string `json:"value"`
}

func Encoder(prop *model.Property) func(w http.ResponseWriter, r *http.Request) interface{} {
	switch prop.ValueType.Type {

	case "number":
		return func(w http.ResponseWriter, r *http.Request) interface{} {
			return WotNumber{store[r.RequestURI].(float64)}
		}
	}

	return nil
}

func Decoder(prop *model.Property) func(w http.ResponseWriter, r *http.Request) interface{} {
	switch prop.ValueType.Type {

	case "number":
		return func(w http.ResponseWriter, r *http.Request) interface{} {
			var v WotNumber
			json.NewDecoder(r.Body).Decode(&v)
			return v.Value
		}
	}

	return nil
}
