package wot

import (
	"reflect"

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
type Server interface {
	Name() string

	InvokeAction(actionName string, parameter interface{}) interface{}

	SetProperty(propertyName string, newValue interface{}) interface{}

	GetProperty(propertyName string) interface{}

	EmitEvent(eventName string, payload interface{}) interface{}

	AddEvent(eventName string, payloadType reflect.Type) *Server

	AddAction(actionName string, inputType interface{}, outputType interface{}) *Server

	AddProperty(propertyName string, contentType reflect.Type) *Server

	OnInvokeAction(actionName string, callback func(interface{}) interface{}) *Server

	OnUpdateProperty(propertyName string, callback func(interface{}) interface{}) *Server

	AddListener(eventName string, listener EventListener) *Server

	RemoveListener(eventName string, listener EventListener) *Server

	RemoveAllListeners(eventName string) *Server

	GetDescription() model.ThingDescription
}
