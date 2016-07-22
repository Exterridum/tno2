package codec

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"

	"github.com/conas/tno2/util/str"
)

// ----- CODEC TYPES

type Encoding string

const (
	ENCODING_JSON Encoding = "JSON"
)

type Codec interface {
	Info() Encoding
	Marshal(io.Writer, interface{}) error
	Unmarshal(io.Reader, interface{}) error
}

// ----- JSON CODEC

type JsonCodec struct{}

func NewJsonCodec() *JsonCodec {
	return &JsonCodec{}
}

func (c *JsonCodec) Info() Encoding {
	return ENCODING_JSON
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
		return errors.New(str.Concat("Error unmarshaling input using ", c.Info(), " codec."))
	}

	return nil
}

// ----- COMMON

func ReaderToByteArray(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
