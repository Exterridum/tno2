package wot

import (
	"log"

	"github.com/conas/tno2/util/sync"
	"github.com/conas/tno2/wot/model"
)

type Server struct {
	pubCh chan<- interface{}
	td    *model.ThingDescription

	propGetCB map[string]func() interface{}
	propSetCB map[string]func(interface{})
	actionCB  map[string]ActionHandler
	eventsCB  map[string]func(interface{})
}

type Device interface {
	Init(initParams map[string]interface{}, s *Server)
}

type ActionHandler func(interface{}, sync.StatusHandler)

type Status int

const (
	OK Status = iota
	UNKNOWN_PROPERTY
	UNKNOWN_ACTION
)

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
		actionCB:  make(map[string]ActionHandler),
		eventsCB:  make(map[string]func(interface{})),
	}
}

func (s *Server) Connect(d Device, initParams map[string]interface{}) {
	d.Init(initParams, s)
}

//FIXME: Create model metadata with map
func (s *Server) propertyExists(name string) bool {
	for _, p := range s.GetDescription().Properties {
		if p.Name == name {
			return true
		}
	}

	return false
}

//FIXME: Create model metadata with map
func (s *Server) actionExists(name string) bool {
	for _, a := range s.GetDescription().Actions {
		if a.Name == name {
			return true
		}
	}

	return false
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
	log.Print("Server -> ", s.GetDescription().Name, " OnUpdateProperty propertyName: ", propertyName)
	s.propSetCB[propertyName] = propUpdateListener
	return s
}

func (s *Server) OnGetProperty(propertyName string, propertyRetriever func() interface{}) *Server {
	log.Print("Server -> ", s.GetDescription().Name, " OnGetProperty propertyName: ", propertyName)
	s.propGetCB[propertyName] = propertyRetriever
	return s
}

func (s *Server) GetProperty(propertyName string) (*sync.Promise, Status) {
	cb, ok := s.propGetCB[propertyName]

	if ok {
		return sync.Async(cb), OK
	} else {
		return nil, UNKNOWN_PROPERTY
	}
}

func (s *Server) SetProperty(propertyName string, newValue interface{}) (*sync.Promise, Status) {
	cb, ok := s.propSetCB[propertyName]

	if ok {
		callable := func() interface{} {
			cb(newValue)
			return nil
		}

		return sync.Async(callable), OK
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
	actionHandler ActionHandler) *Server {
	log.Print("Server -> ", s.GetDescription().Name, " OnInvokeAction actionName: ", actionName)

	s.actionCB[actionName] = actionHandler
	return s
}

func (s *Server) InvokeAction(
	actionName string,
	arg interface{},
	statusHandler sync.StatusHandler) (*sync.StatusPromise, Status) {

	actionHandler, ok := s.actionCB[actionName]

	if ok {
		callable := func(status sync.StatusHandler) interface{} {
			status.Schedule(arg)
			actionHandler(arg, statusHandler)
			return nil
		}

		return sync.AsyncStatus(callable, statusHandler), OK
	} else {
		return nil, UNKNOWN_ACTION
	}
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

func (s *Server) EmitEvent(eventName string, payload interface{}) *sync.Promise {
	return sync.Async(func() interface{} {
		s.eventsCB[eventName](payload)
		return nil
	})
}
