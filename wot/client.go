package wot

import "github.com/conas/tno2/wot/model"

// https://github.com/w3c/wot/tree/master/proposals/restructured-scripting-api#consumedthing
//
// WebIDL
// interface ConsumedThing {
//     readonly attribute DOMString name;
//     Promise<any>  invokeAction(DOMString actionName, any parameter);
//     Promise<any>  setProperty(DOMString propertyName, any newValue);
//     Promise<any>  getProperty(DOMString propertyName);
//     ConsumedThing addListener(DOMString eventName, ThingEventListener listener);
//     ConsumedThing removeListener(DOMString eventName,
//                                  ThingEventListener listener);
//     ConsumedThing removeAllListeners(DOMString eventName);
//     object        getDescription();
// };
type Client interface {
	Name() string

	InvokeAction(actionName string, parameter interface{}) interface{}

	SetProperty(propertyName string, newValue interface{}) interface{}

	GetProperty(propertyName string) interface{}

	AddListener(eventName string, listener func(interface{})) *Client

	RemoveListener(eventName string, listener func(interface{})) *Client

	RemoveAllListeners(eventName string) *Client

	GetDescription() model.ThingDescription
}

func Discover(discoveryType string, filter interface{}) []*Client {
	return nil
}

func ConsumeDescription(td model.ThingDescription) *Client {
	return nil
}

func ConsumeDescriptionUri(uri string) *Client {
	return nil
}
