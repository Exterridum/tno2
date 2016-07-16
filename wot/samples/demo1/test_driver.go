package main

import (
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/model"
)

type TestDriver struct{}

func NewTestDriver() *TestDriver {
	return &TestDriver{}
}

func (d *TestDriver) Init(initParams map[string]interface{}, s *wot.Server) {
	d.addPropsHandler(s.GetDescription())
}

func (d *TestDriver) addPropsHandler(td *model.ThingDescription) {

}
