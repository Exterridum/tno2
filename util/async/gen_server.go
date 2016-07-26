package async

import (
	"sync"
	"time"
)

// Performance testing shows that channels are very slow, need to reimplement using another pattern
type GenServer struct {
	handlers map[MessageType]func(interface{}) interface{}
	in       chan<- interface{}
	out      chan<- interface{}
	actor    *Actor
}

func NewGenServer() *GenServer {
	gs := &GenServer{
		handlers: make(map[MessageType]func(interface{}) interface{}),
	}

	return gs
}

type MessageType int

type Message struct {
	msgType MessageType
	out     chan<- interface{}
	data    interface{}
}

func (gs *GenServer) Handle(msgType MessageType, handler func(interface{}) interface{}) *GenServer {
	gs.handlers[msgType] = handler

	return gs
}

func (gs *GenServer) Start() {
	gs.actor = Spawn(gs.processor, gs.panicHandler)
	gs.in = gs.actor.Channel()

	startShadow()
}

// Seems like golang is unregistering gorutine while recovering.
// During recovery we need to write error back to out channel to report error to client.
// Since Golang is not registering gorutines it reports deadlock when trying to write or read from channel.
// Therefore we need to start shadow gorutin.
// Even solution is working it is not very elegant. Need to find a better one.
func startShadow() {
	shadowMux.Lock()
	defer shadowMux.Unlock()

	if started {
		return
	}

	go func() {
		for {
			time.Sleep(time.Hour * 12)
		}
	}()

	started = true
}

var shadowMux = &sync.Mutex{}
var started = false

func (gs *GenServer) panicHandler(err interface{}) {
	gs.out <- err
}

func (gs *GenServer) Call(msgType MessageType, data interface{}) <-chan interface{} {
	out := make(chan interface{}, 1)

	gs.in <- &Message{
		msgType: msgType,
		out:     out,
		data:    data,
	}

	return out
}

func (gs *GenServer) processor(in <-chan interface{}) {
	for {
		mail := <-in
		msg := mail.(*Message)
		gs.out = msg.out
		gs.out <- gs.handlers[msg.msgType](msg.data)
		close(gs.out)
		gs.out = nil
	}
}
