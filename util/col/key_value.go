package col

import "sync"

type Map struct {
	l *sync.RWMutex
	v map[string]interface{}
}

func NewConcurentMap() *Map {
	return &Map{
		l: &sync.RWMutex{},
		v: make(map[string]interface{}),
	}
}

func AsConcurentMap(m map[string]interface{}) *Map {
	return &Map{
		l: &sync.RWMutex{},
		v: m,
	}
}

func (am *Map) Add(k string, v interface{}) {
	am.l.Lock()
	am.v[k] = v
	am.l.Unlock()
}

func (am *Map) Get(k string) (interface{}, bool) {
	am.l.RLock()
	v, ok := am.v[k]
	am.l.RUnlock()

	return v, ok
}

func (am *Map) Del(k string) {
	am.l.Lock()
	delete(am.v, k)
	am.l.Unlock()
}

func (am *Map) Keys() []string {
	am.l.RLock()

	keys := make([]string, 0, len(am.v))
	for k := range am.v {
		keys = append(keys, k)
	}

	am.l.RUnlock()
	return keys
}

type KeyValue struct {
	K string
	V interface{}
}

func KV(k string, v interface{}) *KeyValue {
	return &KeyValue{
		K: k,
		V: v,
	}
}

func AsMap(kvs ...*KeyValue) map[string]interface{} {
	params := make(map[string]interface{})

	for _, kv := range kvs {
		params[kv.K] = kv.V
	}

	return params
}
