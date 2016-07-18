package main

import (
	"log"
	"time"

	"github.com/conas/tno2/util/sync"
	"github.com/conas/tno2/wot"
)

type TestDevice struct {
	datastore map[string]interface{}
}

func NewTestDevice() *TestDevice {
	return &TestDevice{
		datastore: make(map[string]interface{}),
	}
}

func (d *TestDevice) Init(initParams map[string]interface{}, s *wot.Server) {
	log.Println("TestDriver -> initializing server ->", s.GetDescription().Name)
	d.addPropsHandlers(s)
	d.addActionsHandlers(s)
}

func (d *TestDevice) addPropsHandlers(s *wot.Server) {
	for _, p := range s.GetDescription().Properties {
		log.Print("TestDevice -> found property: ", p.Name, ", writable:", p.Writable)
		s.OnGetProperty(p.Name, d.getPropertyHandler(p.Name))

		if p.Writable {
			s.OnUpdateProperty(p.Name, d.propUpdateListener(p.Name))
		}
	}
}

func (d *TestDevice) getPropertyHandler(name string) func() interface{} {
	return func() interface{} {
		return d.datastore[name]
	}
}

func (d *TestDevice) propUpdateListener(name string) func(newValue interface{}) {
	return func(newValue interface{}) {
		d.datastore[name] = newValue
	}
}

func (d *TestDevice) addActionsHandlers(s *wot.Server) {
	for _, a := range s.GetDescription().Actions {
		s.OnInvokeAction(a.Name, d.longRunningAction(a.Name))
	}
}

type ActionStatus struct {
	Status int
}

func (d *TestDevice) longRunningAction(name string) wot.ActionHandler {
	return func(arg interface{}, status sync.StatusHandler) {
		for i := 0; i < 10; i++ {
			status.Update(&ActionStatus{Status: i})
			time.Sleep(time.Second * 2)
		}

		status.Done(&ActionStatus{Status: 10})
	}
}
