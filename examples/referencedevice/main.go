package main

import (
	"time"

	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	wotServer := server.CreateFromDescriptionUri(refModel)
	device := NewTestDevice()
	wotServer.Connect(device, make(map[string]interface{}))

	http := server.NewHttp(8080)
	http.Bind("/reference-model", wotServer)

	startEventGenerator(wotServer)

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
