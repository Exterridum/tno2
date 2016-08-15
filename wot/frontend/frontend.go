package frontend

import (
	"errors"
	"io"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/server"
)

type Factory func(map[string]interface{}) Frontend

type Frontend interface {
	Bind(ctxPath string, s *server.WotServer)
	Start()
}

// ----- CODEC TYPES

const (
	ENCODING_JSON string = "JSON"
)

type Encoder interface {
	Info() string
	Encode(io.Writer, interface{}) error
	Decode(io.Reader, interface{}) error
}

type EncoderRegistry struct {
	reg *col.Map
}

func NewEncoderRegistry() *EncoderRegistry {
	return &EncoderRegistry{
		reg: col.NewConcurentMap(),
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

	return nil, errors.New(str.Concat("Unsupported frontend encoding: ", code))
}

func (es *EncoderRegistry) Registered() []string {
	return es.reg.Keys()
}
