package server

import (
	"log"

	"github.com/conas/tno2/util/sync"
	"github.com/conas/tno2/wot/model"
)

type WotServer struct {
	pubCh chan<- interface{}
	td    *model.ThingDescription

	propGetCB map[string]func() interface{}
	propSetCB map[string]func(interface{})
	actionCB  map[string]ActionHandler
	eventsCB  map[string]func(interface{})
}

type Device interface {
	Init(initParams map[string]interface{}, s *WotServer)
}

type ActionHandler func(interface{}, sync.StatusHandler)

type Status int

const (
	OK Status = iota
	UNKNOWN_PROPERTY
	UNKNOWN_ACTION
)

func CreateThing(name string) *WotServer {
	return nil
}

func CreateFromDescriptionUri(uri string) *WotServer {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *WotServer {
	return &WotServer{
		td:        td,
		pubCh:     make(chan interface{}),
		propGetCB: make(map[string]func() interface{}),
		propSetCB: make(map[string]func(interface{})),
		actionCB:  make(map[string]ActionHandler),
		eventsCB:  make(map[string]func(interface{})),
	}
}

func (s *WotServer) Connect(d Device, initParams map[string]interface{}) {
	d.Init(initParams, s)
}

//FIXME: Create model metadata with map
func (s *WotServer) propertyExists(name string) bool {
	for _, p := range s.GetDescription().Properties {
		if p.Name == name {
			return true
		}
	}

	return false
}

//FIXME: Create model metadata with map
func (s *WotServer) actionExists(name string) bool {
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
func (s *WotServer) Name() string {
	return s.td.Name
}

func (s *WotServer) GetDescription() *model.ThingDescription {
	return s.td
}

// ----- PROPERTIES HANDLING

func (s *WotServer) AddProperty(propertyName string, property interface{}) *WotServer {
	//Should we update TD
	panic("Add property not implemented!")
}

func (s *WotServer) OnUpdateProperty(propertyName string, propUpdateListener func(newValue interface{})) *WotServer {
	log.Print("Server -> ", s.GetDescription().Name, " OnUpdateProperty propertyName: ", propertyName)
	s.propSetCB[propertyName] = propUpdateListener
	return s
}

func (s *WotServer) OnGetProperty(propertyName string, propertyRetriever func() interface{}) *WotServer {
	log.Print("Server -> ", s.GetDescription().Name, " OnGetProperty propertyName: ", propertyName)
	s.propGetCB[propertyName] = propertyRetriever
	return s
}

func (s *WotServer) GetProperty(propertyName string) (*sync.Promise, Status) {
	cb, ok := s.propGetCB[propertyName]

	if ok {
		return sync.Async(cb), OK
	} else {
		return nil, UNKNOWN_PROPERTY
	}
}

func (s *WotServer) SetProperty(propertyName string, newValue interface{}) (*sync.Promise, Status) {
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

func (s *WotServer) AddAction(actionName string, inputType interface{}, outputType interface{}) *WotServer {
	panic("Add action not implemented!")
}

func (s *WotServer) OnInvokeAction(
	actionName string,
	actionHandler ActionHandler) *WotServer {
	log.Print("Server -> ", s.GetDescription().Name, " OnInvokeAction actionName: ", actionName)

	s.actionCB[actionName] = actionHandler
	return s
}

func (s *WotServer) InvokeAction(
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

func (s *WotServer) AddEvent(eventName string, payloadType interface{}) *WotServer {
	panic("Add event not implemented!")
}

func (s *WotServer) AddListener(eventName string, listener func(interface{})) *WotServer {
	s.eventsCB[eventName] = listener
	return s
}

func (s *WotServer) RemoveListener(eventName string, listener func(interface{})) *WotServer {
	delete(s.eventsCB, eventName)
	return s
}

func (s *WotServer) RemoveAllListeners(eventName string) *WotServer {
	delete(s.eventsCB, eventName)
	return s
}

func (s *WotServer) EmitEvent(eventName string, payload interface{}) *sync.Promise {
	return sync.Async(func() interface{} {
		s.eventsCB[eventName](payload)
		return nil
	})
}
