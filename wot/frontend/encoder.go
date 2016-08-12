package frontend

import (
	"errors"
	"io"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/str"
)

// ----- CODEC TYPES

const (
	ENCODING_JSON string = "JSON"
)

type Encoder interface {
	Info() string
	Marshal(io.Writer, interface{}) error
	Unmarshal(io.Reader, interface{}) error
}

type EncoderRegistry struct {
	reg *async.Map
}

func NewEncoderRegistry() *EncoderRegistry {
	return &EncoderRegistry{
		reg: async.NewConcurentMap(),
	}
}

var Encoders = NewEncoderRegistry()

func (es *EncoderRegistry) Register(e Encoder) {
	es.reg.Add(e.Info(), e)
}

func (es *EncoderRegistry) Get(code string) (Encoder, error) {
	e, ok := es.reg.Get(code)

	if ok {
		return e.(Encoder), nil
	}

	return nil, errors.New(str.Concat("Unsupported encoding: ", code))
}

func (es *EncoderRegistry) Registered() []string {
	return es.reg.Keys()
}
