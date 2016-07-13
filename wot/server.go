package wot

import (
	"reflect"

	"github.com/conas/tno2/concurrent"
	"github.com/conas/tno2/wot/driver"
	"github.com/conas/tno2/wot/model"
)

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
type Server struct {
	td *model.ThingDescription

	invoke chan interface{}

	notification chan interface{}

	actions map[string]func(interface{}) interface{}

	properties map[string]interface{}

	events map[string]reflect.Type
}

func (s *Server) Name() string {
	return s.td.Name
}

func (s *Server) InvokeAction(actionName string, parameter interface{}) *concurent.Promise {
	c := make(chan interface{})
	p := concurent.NewPromise(c)

	s.send(&driver.InvokeActionRQ{
		c,
		actionName,
		parameter,
	})

	return p
}

func (s *Server) send(message *driver.InvokeActionRQ) {
	go func() {
		s.invoke <- message
	}()
}

func (s *Server) SetProperty(propertyName string, newValue interface{}) interface{} {
	return nil
}

func (s *Server) GetProperty(propertyName string) interface{} {
	return nil
}

func (s *Server) EmitEvent(eventName string, payload interface{}) interface{} {
	return nil
}

func (s *Server) AddEvent(eventName string, payloadType reflect.Type) *Server {
	return nil
}

func (s *Server) AddAction(actionName string, inputType interface{}, outputType interface{}) *Server {
	return nil
}

func (s *Server) AddProperty(propertyName string, contentType reflect.Type) *Server {
	return nil
}

func (s *Server) OnInvokeAction(actionName string, callback func(interface{}) interface{}) *Server {
	return nil
}

func (s *Server) OnUpdateProperty(propertyName string, callback func(interface{}) interface{}) *Server {
	return nil
}

func (s *Server) AddListener(eventName string, listener EventListener) *Server {
	return nil
}

func (s *Server) RemoveListener(eventName string, listener EventListener) *Server {
	return nil
}

func (s *Server) RemoveAllListeners(eventName string) *Server {
	return nil
}

func (s *Server) GetDescription() *model.ThingDescription {
	return s.td
}

func CreateThing(name string) *Server {
	return nil
}

func CreateFromDescriptionUri(uri string) *Server {
	return CreateFromDescription(model.Create(uri))
}

func CreateFromDescription(td *model.ThingDescription) *Server {
	return &Server{
		td,
		make(chan interface{}),
		make(chan interface{}),
		make(map[string]func(interface{}) interface{}),
		make(map[string]interface{}),
		make(map[string]reflect.Type),
	}
}
