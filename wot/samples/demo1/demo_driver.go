package main

import (
	"log"

	"github.com/conas/tno2/util/strings"
	"github.com/conas/tno2/wot/driver"
)

type DemoDriver struct {
	actions     map[string]func(arg interface{}) interface{}
	properties  map[string]interface{}
	eventEmiter driver.EventEmiter
}

func NewDemoDriver() *DemoDriver {
	d := DemoDriver{
		actions:    make(map[string]func(arg interface{}) interface{}),
		properties: make(map[string]interface{}),
	}

	d.actions["hello"] = func(arg interface{}) interface{} {
		argMap := arg.(map[string]interface{})
		return strings.Concat("Hello ", argMap["name"].(string), "!")
	}

	return &d
}

func (d *DemoDriver) Init(initParams map[string]interface{}, ee driver.EventEmiter) {
	d.eventEmiter = ee
}

func (d *DemoDriver) InvokeAction(msg *driver.InvokeActionRQ) interface{} {
	return d.actions[msg.ActionName](msg.Parameter)
}

func (d *DemoDriver) GetProperty(msg *driver.GetPropertyRQ) interface{} {
	log.Printf("GetProperty: request -> %v", msg)
	value, ok := d.properties[msg.PropertyName]

	if ok {
		return value
	} else {
		return -1000.0
	}
}

func (d *DemoDriver) SetProperty(msg *driver.SetPropertyRQ) {
	log.Printf("SetProperty: request: %v", msg)
	d.properties[msg.PropertyName] = msg.Value
}
