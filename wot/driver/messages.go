package driver

type Message interface {
	SetChannel(ch chan<- interface{})
	GetChannel() chan<- interface{}
}

type InvokeActionRQ struct {
	ActionName string
	Parameter  interface{}
	callbackCh chan<- interface{}
}

func (m *InvokeActionRQ) SetChannel(ch chan<- interface{}) {
	m.callbackCh = ch
}

func (m *InvokeActionRQ) GetChannel() chan<- interface{} {
	return m.callbackCh
}

type GetPropertyRQ struct {
	callbackCh   chan<- interface{}
	PropertyName string
}

func (m *GetPropertyRQ) SetChannel(ch chan<- interface{}) {
	m.callbackCh = ch
}

func (m *GetPropertyRQ) GetChannel() chan<- interface{} {
	return m.callbackCh
}

type SetPropertyRQ struct {
	callbackCh   chan<- interface{}
	PropertyName string
	Value        interface{}
}

func (m *SetPropertyRQ) SetChannel(ch chan<- interface{}) {
	m.callbackCh = ch
}

func (m *SetPropertyRQ) GetChannel() chan<- interface{} {
	return m.callbackCh
}

type Event struct {
	EventName string
	Payload   interface{}
}
