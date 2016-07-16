package main

import (
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot"
)

type TestDriver struct {
	actions    map[string]func(arg interface{}) interface{}
	properties map[string]interface{}
}

func NewTestDriver() *TestDriver {
	d := TestDriver{
		actions:    make(map[string]func(arg interface{}) interface{}),
		properties: make(map[string]interface{}),
	}

	d.actions["hello"] = func(arg interface{}) interface{} {
		argMap := arg.(map[string]interface{})
		return str.Concat("Hello ", argMap["name"].(string), "!")
	}

	return &d
}

func (d *TestDriver) Init(initParams map[string]interface{}, s *wot.Server) {
}
