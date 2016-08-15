package main

import (
	"time"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/platform"
	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	platform.NewPlatform().
		AddFrontend("http-1", "HTTP", col.KV("port", 8080)).
		AddBackend("mqtt-1", "MQTT", col.KV("url", "tcp://localhost:1883")).
		AddWotServer("reference-dev-1", "file://reference-model.json", "/reference-dev-1", "SIMPLE_URL_CODEC", "mqtt-1", []string{"http-1"}).
		Start()

	// ----- OLD TYPE CONFIG
	// wotServer := server.CreateFromDescriptionUri(refModel)
	// startEventGenerator(wotServer)
	// SimulatorBackend().Bind(wotServer)
	// backend.NewMQTT("tcp://localhost:1883").Bind("/topic", wotServer, &backend.SimpleCodec{})
	// frontend.NewHTTP(8080).Bind("/reference-model", wotServer).Start()
}

func startEventGenerator(wotServet *server.WotServer) {
	go func() {
		for {
			wotServet.EmitEvent("critical-condition-event", "Event Data")
			time.Sleep(time.Second * 5)
		}
	}()
}
