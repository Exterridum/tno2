package wot

import (
	"log"

	"github.com/conas/tno2/util/concurrent"
	"github.com/conas/tno2/wot/model"
)

//TODO: So far only one listener is supported per event

type Server struct {
	pubCh chan<- interface{}
	td    *model.ThingDescription

	propGetCB map[string]func() interface{}
	propSetCB map[string]func(interface{})
	actionCB  map[string]func(interface{}, concurent.StatusHandler)
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
		actionCB:  make(map[string]func(interface{}, concurent.StatusHandler)),
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

type RETURN_CODES int

const (
	OK RETURN_CODES = iota
	UNKNOWN_PROPERTY
)

func (s *Server) AddProperty(propertyName string, property interface{}) *Server {
	//Should we update TD
	panic("Add property not implemented!")
}

func (s *Server) OnUpdateProperty(propertyName string, propUpdateListener func(newValue interface{})) *Server {
	log.Print("Server -> ", s.GetDescription().Name, " OnUpdateProperty propertyName: ", propertyName)
	s.propSetCB[propertyName] = propUpdateListener
	return s
}

func (s *Server) OnGetProperty(propertyName string, propertyRetriever func() interface{}) *Server {
	log.Print("Server -> ", s.GetDescription().Name, " OnGetProperty propertyName: ", propertyName)
	s.propGetCB[propertyName] = propertyRetriever
	return s
}

func (s *Server) GetProperty(propertyName string) (*concurent.Promise, RETURN_CODES) {
	cb, ok := s.propGetCB[propertyName]

	if ok {
		return concurent.Async(cb), OK
	} else {
		return nil, UNKNOWN_PROPERTY
	}
}

func (s *Server) SetProperty(propertyName string, newValue interface{}) (*concurent.Promise, RETURN_CODES) {
	cb, ok := s.propSetCB[propertyName]

	if ok {
		return concurent.Async(func() interface{} {
			cb(newValue)
			return nil
		}), OK
	} else {
		return nil, UNKNOWN_PROPERTY
	}
}

// ----- ACTIONS HANDLING

func (s *Server) AddAction(actionName string, inputType interface{}, outputType interface{}) *Server {
	panic("Add action not implemented!")
}

func (s *Server) OnInvokeAction(
	actionName string,
	actionHandler func(interface{}, concurent.StatusHandler)) *Server {

	s.actionCB[actionName] = actionHandler
	return s
}

func (s *Server) InvokeAction(
	actionName string,
	arg interface{},
	statusHandler concurent.StatusHandler) *concurent.StatusPromise {

	actionHandler := s.actionCB[actionName]

	return concurent.AsyncStatus(
		func(*concurent.StatusHandler) interface{} {
			actionHandler(arg, statusHandler)
			return nil
		},
		statusHandler)
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
