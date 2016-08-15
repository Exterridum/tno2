package async

import log "github.com/Sirupsen/logrus"

// Performance testing shows that channels are slow
type GenServer struct {
	handlers map[MessageType]func(interface{}) interface{}
	in       chan<- interface{}
	prom     *Promise
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
	prom    *Promise
	data    interface{}
}

func (gs *GenServer) HandleCall(msgType MessageType, handler func(interface{}) interface{}) *GenServer {
	gs.handlers[msgType] = handler
	return gs
}

func (gs *GenServer) Start() {
	gs.actor = Spawn(16, gs.processor, gs.panicHandler)
	gs.in = gs.actor.Channel()
}

func (gs *GenServer) panicHandler(err interface{}) {
	log.Info(err)
	gs.prom.Set(err)
}

func (gs *GenServer) Call(msgType MessageType, data interface{}) *Promise {
	prom := NewPromise()

	gs.in <- &Message{
		msgType: msgType,
		prom:    prom,
		data:    data,
	}

	return prom
}

func (gs *GenServer) processor(in <-chan interface{}) {
	for {
		mail := <-in
		msg := mail.(*Message)

		//current promise needs to be cached so in case of panic,
		//panic handler can fulfill the promise
		gs.prom = msg.prom
		r := gs.handlers[msg.msgType](msg.data)
		gs.prom.Set(r)
	}
}
