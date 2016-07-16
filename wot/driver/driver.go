package driver

import "github.com/conas/tno2/util/concurrent"

type Driver interface {
	Init(initParams map[string]interface{}, eventEmiter EventEmiter)
	InvokeAction(msg *InvokeActionRQ) interface{}
	GetProperty(msg *GetPropertyRQ) interface{}
	SetProperty(msg *SetPropertyRQ)
}

type EventEmiter interface {
	EmitEvent(eventName string, payload interface{}) *concurent.Promise
}

type MessageType int

const (
	IARQ MessageType = iota
	GPRQ
	SPRQ
)

type Message interface {
	GetMessageType() MessageType
}

type InvokeActionRQ struct {
	ActionName string
	Parameter  interface{}
}

func (m *InvokeActionRQ) GetMessageType() MessageType {
	return IARQ
}

type GetPropertyRQ struct {
	PropertyName string
}

func (m *GetPropertyRQ) GetMessageType() MessageType {
	return GPRQ
}

type SetPropertyRQ struct {
	PropertyName string
	Value        interface{}
}

func (m *SetPropertyRQ) GetMessageType() MessageType {
	return SPRQ
}

//TODO: Do we need Event?
type Event struct {
	EventName string
	Payload   interface{}
}
