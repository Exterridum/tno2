package main

import "github.com/conas/tno2/wot"

type TestDriver struct {
	datastore map[string]interface{}
}

func NewTestDriver() *TestDriver {
	return &TestDriver{
		datastore: make(map[string]interface{}),
	}
}

func (d *TestDriver) Init(initParams map[string]interface{}, s *wot.Server) {
	d.addPropsHandlers(s)
	d.addActionsHandlers(s)
}

func (d *TestDriver) addPropsHandlers(s *wot.Server) {
	for _, p := range s.GetDescription().Properties {
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

func (d *TestDriver) longRunningAction(name string) func(arg interface{}) {
	return func(arg interface{}) {

	}
}
