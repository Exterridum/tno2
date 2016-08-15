package backend

import (
	"errors"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/server"
)

type Factory func(map[string]interface{}) Backend

type Backend interface {
	Bind(s *server.WotServer, ctxPath string, encoder Encoder)
	Start()
}

const (
	BE_ACTION_RQ        int8 = 0
	BE_ACTION_RS        int8 = 1
	BE_GET_PROP_RQ      int8 = 2
	BE_GET_PROP_RS      int8 = 3
	BE_SET_PROP_RQ      int8 = 4
	BE_SET_PROP_RS      int8 = 5
	BE_EVENT            int8 = 6
	BE_UNKNOWN_MSG_TYPE int8 = 7
)

type Encoder interface {
	Info() string
	Decode(buf []byte) (msgType int8, conversationID string, msgName string, data interface{})
	Encode(msgType int8, conversationID string, msgName string, data interface{}) []byte
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

	return nil, errors.New(str.Concat("Unsupported backend encoding: ", code))
}

func (es *EncoderRegistry) Registered() []string {
	return es.reg.Keys()
}
