package main

import (
	"log"
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	wotServer := server.CreateFromDescriptionUri(refModel)
	setupServer(wotServer)

	http := server.NewHttp(8080)
	http.Bind("/reference-model", wotServer)

	// startEventGenerator(wotServer)

	http.Start()
}

func startEventGenerator(wotServet *server.WotServer) {
	go func() {
		for {
			wotServet.EmitEvent("critical-condition-event", "some payload")
			time.Sleep(time.Second * 5)
		}
	}()
}

type ActionStatus struct {
	Status int
}

func setupServer(s *server.WotServer) {
	log.Println("TestDriver -> initializing server ->", s.GetDescription().Name)

	datastore := make(map[string]interface{})

	addPropsHandlers(s, datastore)
	addActionsHandlers(s)
}

func addPropsHandlers(s *server.WotServer, datastore map[string]interface{}) {
	for _, p := range s.GetDescription().Properties {
		log.Print("TestDevice -> found property: ", p.Name, ", writable:", p.Writable)
		s.OnGetProperty(p.Name, getPropertyHandler(p.Name, datastore))

		if p.Writable {
			s.OnUpdateProperty(p.Name, propUpdateListener(p.Name, datastore))
		}
	}
}

func getPropertyHandler(name string, datastore map[string]interface{}) func() interface{} {
	return func() interface{} {
		return datastore[name]
	}
}

func propUpdateListener(name string, datastore map[string]interface{}) func(newValue interface{}) {
	return func(newValue interface{}) {
		datastore[name] = newValue
	}
}

func addActionsHandlers(s *server.WotServer) {
	for _, a := range s.GetDescription().Actions {
		s.OnInvokeAction(a.Name, longRunningAction(a.Name))
	}
}

func longRunningAction(name string) server.ActionHandler {
	return func(arg interface{}, ph async.ProgressHandler) interface{} {
		for i := 0; i < 10; i++ {
			ph.Update(&ActionStatus{Status: i})
			time.Sleep(time.Second * 2)
		}

		return &ActionStatus{Status: 10}
	}
}
