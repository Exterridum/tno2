package main

import (
	"time"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/platform"
	"github.com/conas/tno2/wot/server"
)

var model = "file://../example-model.json"

func main() {
	p := platform.NewPlatform()
	p.AddFrontend("http-1", "HTTP", col.KV("port", 8080))
	p.AddBackend("mqtt-1", "MQTT", col.KV("url", "tcp://localhost:1883"))
	p.AddWotServer("example-dev", model, "/02-mqtt-example", "SIMPLE_URL_CODEC", "mqtt-1", []string{"http-1"})
	wg := p.Start()

	startEventGenerator(p.WotServer("example-dev"))
	wg.Wait()
}

func startEventGenerator(wotServet *server.WotServer) {
	go func() {
		for {
			wotServet.EmitEvent("critical-condition-event", "Event Data")
			time.Sleep(time.Second * 5)
		}
	}()
}
