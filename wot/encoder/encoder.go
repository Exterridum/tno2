package encoder

import (
	"errors"
	"io"
	"sync"

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

type Encoders struct {
	rwLock   *sync.RWMutex
	encoders map[string]Encoder
}

var Registry Encoders = Encoders{
	rwLock:   &sync.RWMutex{},
	encoders: make(map[string]Encoder),
}

func (es Encoders) Register(e Encoder) {
	es.rwLock.Lock()
	defer es.rwLock.Unlock()

	es.encoders[e.Info()] = e
}

func (es Encoders) Get(code string) (Encoder, error) {
	es.rwLock.RLock()
	defer es.rwLock.RUnlock()

	e, ok := es.encoders[code]

	if ok {
		return e, nil
	} else {
		return nil, errors.New(str.Concat("Unsupported encoding: ", code))
	}
}

func (es Encoders) Registered() []string {
	es.rwLock.RLock()
	defer es.rwLock.RUnlock()

	keys := make([]string, 0, len(es.encoders))
	for k := range es.encoders {
		keys = append(keys, k)
	}

	return keys
}
