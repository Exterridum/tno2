package main

import (
	"time"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/backend"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/platform"
	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	p := platform.NewPlatform()

	fe1 := platform.NewComponentConfig("HTTP", col.KV("port", 8080))
	p.AddFronted("http-1", fe1)

	be1 := platform.NewComponentConfig("MQTT", col.KV("url", "tcp://localhost"), col.KV("port", 1883))
	p.AddBackend("mqtt-1", be1)

	simulator := platform.NewWotServer("file://reference-model.json", "mqtt-1", "http-1")
	p.AddWotServer("simulator", simulator)
	p.Start()

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
