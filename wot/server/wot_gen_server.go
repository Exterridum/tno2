package server

import (
	"log"

	"github.com/conas/tno2/util/async"
)

type Status int

const (
	WOT_OK Status = iota
	WOT_UNKNOWN_ACTION
	WOT_NO_ACTION_HANDLER
	WOT_NO_PROPERTY_GET_HANDLER
	WOT_NO_PROPERTY_SET_HANDLER
	WOT_UNKNOWN_PROPERTY
	WOT_UNKNOWN_EVENT
)

const (
	ACTION_HANDLER_CALL async.MessageType = iota
	GET_PROPERTY
	SET_PROPERTY
)

type ActionHandlerCallMsg struct {
	name string
	arg  interface{}
	ph   async.ProgressHandler
}

type GetPropertyMsg struct {
	name string
}

type SetPropertyMsg struct {
	name  string
	value interface{}
}

type ActionHandler func(interface{}, async.ProgressHandler) interface{}

// WotGentServer provides process isolation for device represented by one goroutine
func newGenServer(wc *WotCore) *async.GenServer {

	gs := async.NewGenServer().
		HandleCall(ACTION_HANDLER_CALL, func(arg interface{}) interface{} {
			msg := arg.(*ActionHandlerCallMsg)

			handler, ok := wc.actionCB[msg.name]

			if !ok {
				return WOT_NO_ACTION_HANDLER
			}

			log.Printf("Action start %s", msg.name)

			//Progress handler scheduled status is set at wot_server level.
			result := handler(msg.arg, msg.ph)
			msg.ph.Done(result)

			return WOT_OK
		}).
		HandleCall(GET_PROPERTY, func(arg interface{}) interface{} {
			msg := arg.(*GetPropertyMsg)
			handler, ok := wc.propGetCB[msg.name]

			if !ok {
				return WOT_NO_PROPERTY_GET_HANDLER
			}

			return handler()
		}).
		HandleCall(SET_PROPERTY, func(arg interface{}) interface{} {
			msg := arg.(*SetPropertyMsg)
			handler, ok := wc.propSetCB[msg.name]

			if !ok {
				return WOT_NO_PROPERTY_SET_HANDLER
			}

			handler(msg.value)

			return WOT_OK
		})

	gs.Start()

	return gs
}
