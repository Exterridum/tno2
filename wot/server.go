package wot

import (
	"github.com/conas/tno2/util/concurrent"
	"github.com/conas/tno2/wot/model"
)

//TODO: So far only one listener is supported per event

type Server struct {
	pubCh chan<- interface{}
	td    *model.ThingDescription

	propGetCB map[string]func() interface{}
	propSetCB map[string]func(interface{})
	actionCB  map[string]func(interface{})
	eventsCB  map[string]func(interface{})
}

type Driver interface {
	Init(initParams map[string]interface{}, s *Server)
}

func CreateThing(name string) *Server {
	return nil
}

func CreateFromDescriptionUri(uri string) *Server {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *Server {
	return &Server{
		td:        td,
		pubCh:     make(chan interface{}),
		propGetCB: make(map[string]func() interface{}),
		propSetCB: make(map[string]func(interface{})),
		actionCB:  make(map[string]func(interface{})),
		eventsCB:  make(map[string]func(interface{})),
	}
}

func (s *Server) ConnectSync(d Driver, initParams map[string]interface{}) {
	d.Init(initParams, s)
}

// ----- AS DEFINED BY WEB IDL
// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#exposedthing
//
// WebIDL
// interface ExposedThing {
//     readonly attribute DOMString name;
//     Promise<any> invokeAction(DOMString actionName, any parameter);
//     Promise<any> setProperty(DOMString propertyName, any newValue);
//     Promise<any> getProperty(DOMString propertyName);
//     Promise<any> emitEvent(DOMString eventName, any payload);
//     ExposedThing addEvent(DOMString eventName, object payloadType);
//     ExposedThing addAction(DOMString actionName,
//                            object inputType,
//                            object outputType);
//     ExposedThing addProperty(DOMString propertyName, object contentType);
//     ExposedThing onInvokeAction(DOMString actionName, ActionHandler callback);
//     ExposedThing onUpdateProperty(DOMString propertyName,
//                                   PropertyChangeListener callback);
//     ExposedThing addListener(DOMString eventName, ThingEventListener listener);
//     ExposedThing removeListener(DOMString eventName,
//                                 ThingEventListener listener);
//     ExposedThing removeAllListeners(DOMString eventName);
//     object       getDescription();
// };
func (s *Server) Name() string {
	return s.td.Name
}

func (s *Server) GetDescription() *model.ThingDescription {
	return s.td
}

// ----- PROPERTIES HANDLING

func (s *Server) AddProperty(propertyName string, property interface{}) *Server {
	//Should we update TD
	panic("Add property not implemented!")
}

func (s *Server) OnUpdateProperty(propertyName string, propUpdateListener func(newValue interface{})) *Server {
	s.propSetCB[propertyName] = propUpdateListener
	return s
}

func (s *Server) OnGetProperty(propertyName string, callback func() interface{}) *Server {
	s.propGetCB[propertyName] = callback
	return s
}

func (s *Server) GetProperty(propertyName string) *concurent.Promise {
	return concurent.Async(s.propGetCB[propertyName])
}

func (s *Server) SetProperty(propertyName string, newValue interface{}) *concurent.Promise {
	return concurent.Async(func() interface{} {
		s.propSetCB[propertyName](newValue)
		return nil
	})
}

// ----- ACTIONS HANDLING

func (s *Server) AddAction(actionName string, inputType interface{}, outputType interface{}) *Server {
	panic("Add action not implemented!")
}

func (s *Server) OnInvokeAction(actionName string, actionHandler func(arg interface{})) *Server {
	s.actionCB[actionName] = actionHandler
	return s
}

func (s *Server) InvokeAction(actionName string, arg interface{}) *concurent.Promise {
	return concurent.Async(func() interface{} {
		s.actionCB[actionName](arg)
		return nil
	})
}

// ----- EVENTS HANDLING

func (s *Server) AddEvent(eventName string, payloadType interface{}) *Server {
	panic("Add event not implemented!")
}

func (s *Server) AddListener(eventName string, listener func(interface{})) *Server {
	s.eventsCB[eventName] = listener
	return s
}

func (s *Server) RemoveListener(eventName string, listener func(interface{})) *Server {
	delete(s.eventsCB, eventName)
	return s
}

func (s *Server) RemoveAllListeners(eventName string) *Server {
	delete(s.eventsCB, eventName)
	return s
}

func (s *Server) EmitEvent(eventName string, payload interface{}) *concurent.Promise {
	return concurent.Async(func() interface{} {
		s.eventsCB[eventName](payload)
		return nil
	})
}
