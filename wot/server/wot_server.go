package server

import (
	"sync"
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/model"
)

// ----- AS DEFINED BY WEB IDL
// http://w3c.github.io/wot/current-practices/wot-practices.html#idl-def-exposedthing
// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#exposedthing

type WotServer struct {
	td     *model.ThingDescription
	gs     *async.GenServer
	events *eventsMap
}

type EventListener struct {
	ID string
	CB func(interface{})
}

type Event struct {
	Name      string      `json:"name,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func newEvent(eventName string, data interface{}) *Event {
	return &Event{
		Name:      eventName,
		Timestamp: time.Now(),
		Data:      data,
	}
}

type eventsMap struct {
	lock     *sync.RWMutex
	eventsCB map[string][]*EventListener
}

func (em *eventsMap) addEvent(eventName string) Status {
	em.lock.Lock()
	defer em.lock.Unlock()

	_, ok := em.eventsCB[eventName]

	if !ok {
		em.eventsCB[eventName] = make([]*EventListener, 0)
	}

	return WOT_OK
}

func (em *eventsMap) addListener(eventName string, listener *EventListener) Status {
	em.lock.Lock()
	defer em.lock.Unlock()

	_, ok := em.eventsCB[eventName]

	if !ok {
		return WOT_UNKNOWN_EVENT
	}

	em.eventsCB[eventName] = append(em.eventsCB[eventName], listener)

	return WOT_OK
}

func (em *eventsMap) listeners(eventName string) ([]*EventListener, Status) {
	em.lock.RLock()
	defer em.lock.RUnlock()

	listeners, ok := em.eventsCB[eventName]

	if !ok {
		return nil, WOT_UNKNOWN_EVENT
	}

	return listeners, WOT_OK
}

func CreateThing(name string) *WotServer {
	return nil
}

func CreateFromDescriptionUri(uri string) *WotServer {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *WotServer {
	events := &eventsMap{
		lock:     &sync.RWMutex{},
		eventsCB: make(map[string][]*EventListener),
	}

	for _, ev := range td.Events {
		events.addEvent(ev.Name)
	}

	return &WotServer{
		td:     td,
		gs:     setup(),
		events: events,
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
	s.events.addEvent(eventName)
	return s
}

func (s *WotServer) AddListener(eventName string, listener *EventListener) *WotServer {
	//FIXME we should check for event rpesense
	s.events.addListener(eventName, listener)
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

func (s *WotServer) EmitEvent(eventName string, data interface{}) Status {
	listeners, status := s.events.listeners(eventName)

	if status != WOT_OK {
		return status
	}

	// TODO: Check panic safety
	async.Run(func() interface{} {
		event := newEvent(eventName, data)
		for _, eventListener := range listeners {
			eventListener.CB(event)
		}
		return nil
	})

	return WOT_OK
}
