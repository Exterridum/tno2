package main

import "github.com/conas/tno2/wot"

type SimpleDriver struct {
	db map[string]interface{}
}

func NewSimpleDriver() *SimpleDriver {
	return &SimpleDriver{
		db: make(map[string]interface{}),
	}
}

func (d *SimpleDriver) Init(initParams map[string]interface{}, s *wot.Server) {

	s.OnGetProperty("temperature", func() interface{} {
		return d.db["temperature"]
	})

	s.OnUpdateProperty("temperature", func(newValue interface{}) {
		d.db["temperature"] = newValue
	})
}
