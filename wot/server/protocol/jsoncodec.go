package protocol

import (
	"encoding/json"
	"io"

	"github.com/conas/tno2/wot/model"
)

type WotNumber struct {
	Value float64 `json:"value"`
}

type WotString struct {
	Value string `json:"value"`
}

type WotBoolean struct {
	Value bool `json:"value"`
}

type WotObject struct {
	Value interface{} `json:"value"`
}

func Encoder(prop *model.Property) func(value interface{}) interface{} {
	switch prop.ValueType.Type {

	case "number":
		return func(value interface{}) interface{} {
			return WotNumber{value.(float64)}
		}

	case "string":
		return func(value interface{}) interface{} {
			return WotString{value.(string)}
		}

	case "boolean":
		return func(value interface{}) interface{} {
			return WotBoolean{value.(bool)}
		}

	case "object":
		return func(value interface{}) interface{} {
			return WotObject{value.(interface{})}
		}

	default:
		return nil
	}
}

func Decoder(prop *model.Property) func(r io.Reader) interface{} {
	switch prop.ValueType.Type {

	case "number":
		return func(r io.Reader) interface{} {
			var v WotNumber
			json.NewDecoder(r).Decode(&v)
			return v.Value
		}

	case "string":
		return func(r io.Reader) interface{} {
			var v WotString
			json.NewDecoder(r).Decode(&v)
			return v.Value
		}

	case "boolean":
		return func(r io.Reader) interface{} {
			var v WotBoolean
			json.NewDecoder(r).Decode(&v)
			return v.Value
		}

	case "object":
		return func(r io.Reader) interface{} {
			var v WotObject
			json.NewDecoder(r).Decode(&v)
			return v.Value
		}

	default:
		return nil
	}
}
