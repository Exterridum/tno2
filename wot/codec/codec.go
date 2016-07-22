package codec

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// ----- CODEC TYPES

type MediaType string

const (
	MEDIA_TYPE_JSON MediaType = "application/json"
)

type Codec interface {
	Info() MediaType
	Marshal(io.Writer, interface{}) error
	Unmarshal(io.Reader, interface{}) error
}

// ----- JSON CODEC

type JsonCodec struct{}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Info() MediaType {
	return MEDIA_TYPE_JSON
}

func (c *JsonCodec) Marshal(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (c *JsonCodec) Unmarshal(r io.Reader, t interface{}) error {
	data, err := ReaderToByteArray(r)

	if err != nil {
		return err
	}

	err = json.Unmarshal(data, t)

	if err != nil {
		return err
	}

	return err
}

// ----- COMMON

func ReaderToByteArray(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
