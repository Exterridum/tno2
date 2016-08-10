package async

// Performance testing shows that channels are slow
type GenServer struct {
	handlers map[MessageType]func(interface{}) interface{}
	in       chan<- interface{}
	res      *Promise
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
	res     *Promise
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
	gs.res.Set(err)
}

func (gs *GenServer) Call(msgType MessageType, data interface{}) *Promise {
	res := NewPromise()

	gs.in <- &Message{
		msgType: msgType,
		res:     res,
		data:    data,
	}

	return res
}

func (gs *GenServer) processor(in <-chan interface{}) {
	for {
		mail := <-in
		msg := mail.(*Message)
		gs.res = msg.res
		r := gs.handlers[msg.msgType](msg.data)
		gs.res.Set(r)
	}
}
