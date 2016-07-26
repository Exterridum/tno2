package server

import (
	"log"
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/model"
)

// ----- AS DEFINED BY WEB IDL
// http://w3c.github.io/wot/current-practices/wot-practices.html#idl-def-exposedthing
// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#exposedthing

type WotServer struct {
	td        *model.ThingDescription
	propGetCB map[string]func() interface{}
	propSetCB map[string]func(interface{})
	actionCB  map[string]ActionHandler
	eventsCB  map[string][]*EventListener
}

type ActionHandler func(interface{}, async.StatusHandler)

type EventListener struct {
	ID string
	CB func(interface{})
}

type Event struct {
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Event     interface{} `json:"event,omitempty"`
}

type Status int

const (
	WOT_OK Status = iota
	WOT_UNKNOWN_PROPERTY
	WOT_UNKNOWN_ACTION
	WOT_UNKNOWN_EVENT
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
		propGetCB: make(map[string]func() interface{}),
		propSetCB: make(map[string]func(interface{})),
		actionCB:  make(map[string]ActionHandler),
		eventsCB:  make(map[string][]*EventListener),
	}
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

func (s *WotServer) GetProperty(propertyName string) (*async.Promise, Status) {
	cb, ok := s.propGetCB[propertyName]

	if ok {
		return async.Run(cb), WOT_OK
	} else {
		return nil, WOT_UNKNOWN_PROPERTY
	}
}

func (s *WotServer) SetProperty(propertyName string, newValue interface{}) (*async.Promise, Status) {
	cb, ok := s.propSetCB[propertyName]

	if ok {
		callable := func() interface{} {
			cb(newValue)
			return nil
		}

		return async.Run(callable), WOT_OK
	} else {
		return nil, WOT_UNKNOWN_PROPERTY
	}
}

// ----- ACTIONS HANDLING

func (s *WotServer) AddAction(actionName string, inputType interface{}, outputType interface{}) *WotServer {
	panic("Add action not implemented!")
}

func (s *WotServer) OnInvokeAction(
	actionName string,
	actionHandler ActionHandler) *WotServer {

	s.actionCB[actionName] = actionHandler
	return s
}

func (s *WotServer) InvokeAction(
	actionName string,
	arg interface{},
	statusHandler async.StatusHandler) (*async.StatusPromise, Status) {

	actionHandler, ok := s.actionCB[actionName]

	if !ok {
		return nil, WOT_UNKNOWN_ACTION
	}

	callable := func(status async.StatusHandler) interface{} {
		status.Schedule(arg)
		actionHandler(arg, statusHandler)
		return nil
	}

	return async.RunWithStatus(callable, statusHandler), WOT_OK
}

// ----- EVENTS HANDLING

func (s *WotServer) AddEvent(eventName string, payloadType interface{}) *WotServer {
	panic("Add event not implemented!")
}

func (s *WotServer) AddListener(eventName string, listener *EventListener) *WotServer {
	s.eventsCB[eventName] = append(s.eventsCB[eventName], listener)
	return s
}

func (s *WotServer) RemoveListener(eventName string, listener func(interface{})) *WotServer {
	//FIXME: How to identify which listener to remove
	delete(s.eventsCB, eventName)
	return s
}

func (s *WotServer) RemoveAllListeners(eventName string) *WotServer {
	delete(s.eventsCB, eventName)
	return s
}

func (s *WotServer) EmitEvent(eventName string, payload interface{}) (*async.Promise, Status) {
	listeners, ok := s.eventsCB[eventName]

	if ok {
		return async.Run(func() interface{} {
			ev := Event{
				Event:     payload,
				Timestamp: time.Now(),
			}
			for _, eventListener := range listeners {
				eventListener.CB(ev)
			}

			return nil

		}), WOT_OK
	}

	return nil, WOT_OK
}
