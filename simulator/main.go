package main

import (
	"time"

	"github.com/conas/tno2/wot/backend"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	wotServer := server.CreateFromDescriptionUri(refModel)
	// startEventGenerator(wotServer)
	// SimulatorBackend().Bind(wotServer)
	backend.NewMQTT("tcp://localhost:1883").Bind("/topic", wotServer, &backend.SimpleCodec{})
	frontend.NewHTTP(8080).Bind("/reference-model", wotServer).Start()
}

func startEventGenerator(wotServet *server.WotServer) {
	go func() {
		for {
			wotServet.EmitEvent("critical-condition-event", "Event Data")
			time.Sleep(time.Second * 5)
		}
	}()
}
