package server

import (
	"log"
	"time"

	"github.com/conas/tno2/util/async"
)

const (
	WOT_OK Status = iota
	WOT_UNKNOWN_PROPERTY
	WOT_UNKNOWN_ACTION
	WOT_UNKNOWN_EVENT
)

const (
	ACTION_HANDLER_CALL async.MessageType = iota
	ACTION_HANDLER_ADD
	GET_PROPERTY
	GET_PROPERTY_HANDLER
	SET_PROPERTY
	SET_PROPERTY_HANDLER
	EVENT_LISTENER_ADD
	EVENT_EMIT
)

type ActionHandlerAddMsg struct {
	name string
	fn   ActionHandler
}

type ActionHandlerCallMsg struct {
	name string
	arg  interface{}
	sh   async.StatusHandler
}

type GetPropertyMsg struct {
	name string
}

type GetPropertyHandlerMsg struct {
	name string
	fn   func() interface{}
}

type SetPropertyMsg struct {
	name  string
	value interface{}
}

type SetPropertyHandlerMsg struct {
	name string
	fn   func(interface{})
}

type EventEmitMsg struct {
	name string
	data interface{}
}

type EventListenerAddMsg struct {
	name     string
	listener *EventListener
}

type ActionHandler func(interface{}, async.StatusHandler) interface{}

type EventListener struct {
	ID string
	CB func(interface{})
}

type Event struct {
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Event     interface{} `json:"event,omitempty"`
}

type Status int

// WotGentServer provides process isolation for device represented by one goroutine
func setup() *async.GenServer {
	propGetCB := make(map[string]func() interface{})
	propSetCB := make(map[string]func(interface{}))
	actionCB := make(map[string]ActionHandler)
	eventsCB := make(map[string][]*EventListener)

	gs := async.NewGenServer().
		HandleCall(ACTION_HANDLER_ADD, func(arg interface{}) interface{} {
			msg := arg.(*ActionHandlerAddMsg)
			actionCB[msg.name] = msg.fn

			return WOT_OK
		}).
		HandleCall(ACTION_HANDLER_CALL, func(arg interface{}) interface{} {
			msg := arg.(*ActionHandlerCallMsg)
			handler, ok := actionCB[msg.name]

			if !ok {
				return WOT_UNKNOWN_ACTION
			}

			log.Printf("Action start %s", msg.name)

			msg.sh.Schedule(arg)
			result := handler(msg.arg, msg.sh)
			msg.sh.Done(result)

			return WOT_OK
		}).
		HandleCall(GET_PROPERTY, func(arg interface{}) interface{} {
			msg := arg.(*GetPropertyMsg)
			handler, ok := propGetCB[msg.name]

			if !ok {
				return WOT_UNKNOWN_PROPERTY
			}

			return handler()
		}).
		HandleCall(GET_PROPERTY_HANDLER, func(arg interface{}) interface{} {
			msg := arg.(*GetPropertyHandlerMsg)
			propGetCB[msg.name] = msg.fn

			return WOT_OK
		}).
		HandleCall(SET_PROPERTY, func(arg interface{}) interface{} {
			msg := arg.(*SetPropertyMsg)
			handler, ok := propSetCB[msg.name]

			if !ok {
				return WOT_UNKNOWN_PROPERTY
			}

			handler(msg.value)

			return WOT_OK
		}).
		HandleCall(SET_PROPERTY_HANDLER, func(arg interface{}) interface{} {
			msg := arg.(*SetPropertyHandlerMsg)
			propSetCB[msg.name] = msg.fn

			return WOT_OK
		}).
		HandleCall(EVENT_EMIT, func(arg interface{}) interface{} {
			msg := arg.(*EventEmitMsg)

			listeners, ok := eventsCB[msg.name]

			if !ok {
				return WOT_UNKNOWN_EVENT
			}

			// TODO: Check panic safety
			async.Run(func() interface{} {
				for _, eventListener := range listeners {
					eventListener.CB(newEvent(msg.data))
				}
				return nil
			})

			return WOT_OK
		}).
		HandleCall(EVENT_LISTENER_ADD, func(arg interface{}) interface{} {
			msg := arg.(*EventListenerAddMsg)
			eventsCB[msg.name] = append(eventsCB[msg.name], msg.listener)

			return WOT_OK
		})

	gs.Start()

	return gs
}

func newEvent(data interface{}) *Event {
	return &Event{
		Event:     data,
		Timestamp: time.Now(),
	}
}
