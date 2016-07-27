package server

import (
	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/model"
)

// ----- AS DEFINED BY WEB IDL
// http://w3c.github.io/wot/current-practices/wot-practices.html#idl-def-exposedthing
// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#exposedthing

type WotServer struct {
	td *model.ThingDescription
	gs *async.GenServer
}

func CreateThing(name string) *WotServer {
	return nil
}

func CreateFromDescriptionUri(uri string) *WotServer {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *WotServer {
	return &WotServer{
		td: td,
		gs: setup(),
	}
}

func (s *WotServer) Name() string {
	return s.td.Name
}

func (s *WotServer) GetDescription() *model.ThingDescription {
	return s.td
}

func (s *WotServer) AddProperty(propertyName string, property interface{}) *WotServer {
	//Should we update TD?
	panic("Add property not implemented!")
}

func (s *WotServer) OnUpdateProperty(propertyName string, propUpdateListener func(newValue interface{})) *WotServer {
	s.gs.Call(SET_PROPERTY_HANDLER, &SetPropertyHandlerMsg{
		name: propertyName,
		fn:   propUpdateListener,
	})

	return s
}

func (s *WotServer) OnGetProperty(propertyName string, propertyRetriever func() interface{}) *WotServer {
	s.gs.Call(GET_PROPERTY_HANDLER, &GetPropertyHandlerMsg{
		name: propertyName,
		fn:   propertyRetriever,
	})

	return s
}

func (s *WotServer) GetProperty(propertyName string) *async.Value {
	return s.gs.Call(GET_PROPERTY, &GetPropertyMsg{
		name: propertyName,
	})
}

func (s *WotServer) SetProperty(propertyName string, newValue interface{}) *async.Value {
	return s.gs.Call(SET_PROPERTY, &SetPropertyMsg{
		name:  propertyName,
		value: newValue,
	})
}

func (s *WotServer) AddAction(actionName string, inputType interface{}, outputType interface{}) *WotServer {
	panic("Add action not implemented!")
}

func (s *WotServer) OnInvokeAction(actionName string, actionHandler ActionHandler) *WotServer {
	s.gs.Call(ACTION_HANDLER_ADD, &ActionHandlerAddMsg{
		name: actionName,
		fn:   actionHandler,
	})

	return s
}

func (s *WotServer) InvokeAction(actionName string, arg interface{}, ph async.ProgressHandler) *async.Value {
	return s.gs.Call(ACTION_HANDLER_CALL, &ActionHandlerCallMsg{
		name: actionName,
		arg:  arg,
		ph:   ph,
	})
}

func (s *WotServer) AddEvent(eventName string, payloadType interface{}) *WotServer {
	panic("Add event not implemented!")
}

func (s *WotServer) AddListener(eventName string, listener *EventListener) *WotServer {
	s.gs.Call(EVENT_LISTENER_ADD, &EventListenerAddMsg{
		name:     eventName,
		listener: listener,
	})
	return s
}

func (s *WotServer) RemoveListener(eventName string, listener func(interface{})) *WotServer {
	//FIXME: How to identify which listener to remove
	panic("Remove Listener not implemented!")
	return s
}

func (s *WotServer) RemoveAllListeners(eventName string) *WotServer {
	panic("Remove All Listeners not implemented!")
	return s
}

func (s *WotServer) EmitEvent(eventName string, payload interface{}) *async.Value {
	return s.gs.Call(EVENT_EMIT, &EventEmitMsg{
		name: eventName,
		data: payload,
	})
}
