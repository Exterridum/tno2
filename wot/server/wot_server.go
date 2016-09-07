package server

import (
	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/model"
)

// ----- AS DEFINED BY WEB IDL
// http://w3c.github.io/wot/current-practices/wot-practices.html#idl-def-exposedthing
// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#exposedthing

type WotServer struct {
	core *WotCore
	gs   *async.GenServer
}

func CreateThing(name string) *WotServer {
	return nil
}

func CreateFromDescriptionUri(uri string) *WotServer {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *WotServer {
	core := NewWotCoreFromTD(td)
	gs := newGenServer(core)

	return &WotServer{
		core: core,
		gs:   gs,
	}
}

func (s *WotServer) Name() string {
	return s.core.td.Name
}

// ----- DEFINITIONS

func (s *WotServer) AddProperty(propertyName string, property model.Property) *WotServer {
	s.core.PropertyAdd(property)
	return s
}

func (s *WotServer) OnGetProperty(propertyName string, propertyRetriever func() interface{}) *WotServer {
	if s.core.checkProperty(propertyName) == false {
		panic("Property not defined.")
	}
	s.core.propGetCB[propertyName] = propertyRetriever
	return s
}

func (s *WotServer) OnUpdateProperty(propertyName string, propUpdateListener func(newValue interface{})) *WotServer {
	if s.core.checkProperty(propertyName) == false {
		panic("Property not defined.")
	}
	s.core.propSetCB[propertyName] = propUpdateListener
	return s
}

func (s *WotServer) AddAction(actionName string, inputType model.InputData, outputType model.OutputData) *WotServer {
	action := model.Action{
		Name:       actionName,
		InputData:  inputType,
		OutputData: outputType,
	}
	s.core.ActionAdd(action)
	return s
}

func (s *WotServer) OnInvokeAction(actionName string, actionHandler ActionHandler) *WotServer {
	if s.core.checkAction(actionName) == false {
		panic("Action not defined.")
	}
	s.core.actionCB[actionName] = actionHandler
	return s
}

func (s *WotServer) AddEvent(eventName string, event model.Event) *WotServer {
	s.core.EventAdd(event)
	return s
}

func (s *WotServer) AddListener(eventName string, listener *EventListener) *WotServer {
	if s.core.checkEvent(eventName) == false {
		panic("Event not defined.")
	}
	s.core.addListener(eventName, listener)
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

// ----- CALLS

func (s *WotServer) GetDescription() *model.ThingDescription {
	return s.core.td
}

func (s *WotServer) GetProperty(propertyName string) *async.Promise {
	return s.gs.Call(GET_PROPERTY, &GetPropertyMsg{
		name: propertyName,
	})
}

func (s *WotServer) SetProperty(propertyName string, newValue interface{}) *async.Promise {
	return s.gs.Call(SET_PROPERTY, &SetPropertyMsg{
		name:  propertyName,
		value: newValue,
	})
}

func (s *WotServer) InvokeAction(actionName string, arg interface{}, ph async.ProgressHandler) *async.Promise {
	ph.Schedule(arg)

	return s.gs.Call(ACTION_CALL, &ActionHandlerCallMsg{
		name: actionName,
		arg:  arg,
		ph:   ph,
	})
}

func (s *WotServer) EmitEvent(eventName string, data interface{}) Status {
	listeners, status := s.core.listeners(eventName)

	if status != WOT_OK {
		return status
	}

	async.Run(func() interface{} {
		event := newEvent(eventName, data)
		for _, eventListener := range listeners {
			eventListener.CB(event)
		}
		return nil
	})

	return WOT_OK
}
