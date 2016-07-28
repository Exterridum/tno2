package server

import (
	"log"

	"github.com/conas/tno2/util/async"
)

type Status int

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
)

type ActionHandlerAddMsg struct {
	name string
	fn   ActionHandler
}

type ActionHandlerCallMsg struct {
	name string
	arg  interface{}
	ph   async.ProgressHandler
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

type ActionHandler func(interface{}, async.ProgressHandler) interface{}

// WotGentServer provides process isolation for device represented by one goroutine
func setup() *async.GenServer {
	propGetCB := make(map[string]func() interface{})
	propSetCB := make(map[string]func(interface{}))
	actionCB := make(map[string]ActionHandler)

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

			//Progress handler scheduled status is set at transport level.
			result := handler(msg.arg, msg.ph)
			msg.ph.Done(result)

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
		})

	gs.Start()

	return gs
}
