package server

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/conas/tno2/wot/model"
)

type WotCore struct {
	l          *sync.RWMutex
	td         *model.ThingDescription
	properties map[string]model.Property
	actions    map[string]model.Action
	events     map[string]model.Event
	propGetCB  map[string]func() interface{}
	propSetCB  map[string]func(interface{})
	actionCB   map[string]ActionHandler
	eventsCB   map[string][]*EventListener
}

type EventListener struct {
	ID string
	CB func(interface{})
}

type Event struct {
	Event     string      `json:"event,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func newEvent(eventName string, data interface{}) *Event {
	return &Event{
		Event:     eventName,
		Timestamp: time.Now(),
		Data:      data,
	}
}

type eventsListeners struct {
	lock     *sync.RWMutex
	eventsCB map[string][]*EventListener
}

func NewWotCore() *WotCore {
	return &WotCore{
		l:          &sync.RWMutex{},
		properties: make(map[string]model.Property),
		actions:    make(map[string]model.Action),
		events:     make(map[string]model.Event),
		propGetCB:  make(map[string]func() interface{}),
		propSetCB:  make(map[string]func(interface{})),
		actionCB:   make(map[string]ActionHandler),
		eventsCB:   make(map[string][]*EventListener),
	}
}

func NewWotCoreFromTD(td *model.ThingDescription) *WotCore {
	log.Info("Parsing ThingDescription to WotCore")

	wc := NewWotCore()
	wc.td = td

	for _, p := range td.Properties {
		wc.properties[p.Name] = p
	}
	for _, a := range td.Actions {
		wc.actions[a.Name] = a
	}
	for _, e := range td.Events {
		wc.events[e.Name] = e
		wc.eventsCB[e.Name] = make([]*EventListener, 0)
	}

	return wc
}

func (wc *WotCore) PropertyAdd(p model.Property) {
	wc.l.Lock()
	defer wc.l.Unlock()

	wc.td.Properties = append(wc.td.Properties, p)
	wc.properties[p.Name] = p
}

func (wc *WotCore) ActionAdd(a model.Action) {
	wc.l.Lock()
	defer wc.l.Unlock()

	wc.td.Actions = append(wc.td.Actions, a)
	wc.actions[a.Name] = a
}

func (wc *WotCore) EventAdd(e model.Event) {
	wc.l.Lock()
	defer wc.l.Unlock()

	wc.td.Events = append(wc.td.Events, e)
	wc.events[e.Name] = e
	wc.eventsCB[e.Name] = make([]*EventListener, 0)
}

func (wc *WotCore) checkProperty(name string) bool {
	wc.l.RLock()
	defer wc.l.RUnlock()

	_, ok := wc.properties[name]
	return ok
}

func (wc *WotCore) checkAction(name string) bool {
	wc.l.RLock()
	defer wc.l.RUnlock()

	_, ok := wc.actions[name]
	return ok
}

func (wc *WotCore) checkEvent(name string) bool {
	wc.l.RLock()
	defer wc.l.RUnlock()

	_, ok := wc.events[name]
	return ok
}

func (wc *WotCore) addListener(eventName string, listener *EventListener) Status {
	wc.l.Lock()
	defer wc.l.Unlock()

	_, ok := wc.events[eventName]

	if !ok {
		return WOT_UNKNOWN_EVENT
	}

	wc.eventsCB[eventName] = append(wc.eventsCB[eventName], listener)

	return WOT_OK
}

func (wc *WotCore) listeners(eventName string) ([]*EventListener, Status) {
	wc.l.RLock()
	defer wc.l.RUnlock()

	listeners, ok := wc.eventsCB[eventName]

	if !ok {
		return nil, WOT_UNKNOWN_EVENT
	}

	return listeners, WOT_OK
}
