package main

import (
	"log"
	"time"

	"github.com/conas/tno2/util/concurrent"
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot"
)

type TestDriver struct {
	datastore map[string]interface{}
}

func NewTestDriver() *TestDriver {
	return &TestDriver{
		datastore: make(map[string]interface{}),
	}
}

func (d *TestDriver) Init(initParams map[string]interface{}, s *wot.Server) {
	log.Println("TestDriver -> initializing server ->", s.GetDescription().Name)
	d.addPropsHandlers(s)
	d.addActionsHandlers(s)
}

func (d *TestDriver) addPropsHandlers(s *wot.Server) {
	for _, p := range s.GetDescription().Properties {
		log.Print("TestDriver -> found property: ", p.Name, ", writable:", p.Writable)
		s.OnGetProperty(p.Name, d.getPropertyHandler(p.Name))

		if p.Writable {
			s.OnUpdateProperty(p.Name, d.propUpdateListener(p.Name))
		}
	}
}

func (d *TestDriver) getPropertyHandler(name string) func() interface{} {
	return func() interface{} {
		return d.datastore[name]
	}
}

func (d *TestDriver) propUpdateListener(name string) func(newValue interface{}) {
	return func(newValue interface{}) {
		d.datastore[name] = newValue
	}
}

func (d *TestDriver) addActionsHandlers(s *wot.Server) {
	for _, a := range s.GetDescription().Actions {
		s.OnInvokeAction(a.Name, d.longRunningAction(a.Name))
	}
}

func (d *TestDriver) longRunningAction(name string) func(arg interface{}, statusHandler concurent.StatusHandler) {
	return func(arg interface{}, statusHandler concurent.StatusHandler) {
		for i := 0; i < 10; i++ {
			statusHandler(i, str.Concat("Action -> ", name, ", Progress -> ", i*10, "%"))
			time.Sleep(time.Second * 5)
		}

		statusHandler(100, "Action done.")
	}
}
