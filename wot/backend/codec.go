package backend

import "github.com/conas/tno2/util/col"

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

type Codec interface {
	Decode(buf []byte) (msgType int8, conversationID string, msgName string, data interface{})
	Encode(msgType int8, conversationID string, msgName string, data interface{}) []byte
}

type CodecRegistry struct {
	reg *col.Map
}

func NewCodecRegistry() *CodecRegistry {
	return &CodecRegistry{
		reg: col.NewConcurentMap(),
	}
}

var Encoders = NewCodecRegistry()
