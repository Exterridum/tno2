package main

import (
	"fmt"
	"reflect"
)

//https://github.com/nats-io/nats/blob/master/enc.go
type handlers struct {
	cbs map[string]func(interface{})
}

type Handler interface{}

func newHandlers() *handlers {
	return &handlers{
		cbs: make(map[string]func(interface{})),
	}
}

func (*handlers) Subscribe(subject string, cb interface{}) {
	cbv := reflect.ValueOf(cb)
	cbv.Call

	fmt.Println("type %v", reflect.ValueOf(cb))
}

type Person struct{}

func main() {
	h := newHandlers()

	h.Subscribe("name", func(foo string) {

	})

	h.Subscribe("name", func(person Person) {

	})
}
