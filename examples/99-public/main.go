package main

import (
	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/platform"
)

var dhtModel = "file://dht-model.json"

func main() {
	p := platform.NewPlatform()
	p.AddFrontend("http-1", "HTTP", col.KV("port", 8080))
	p.AddBackend("mqtt-1", "MQTT-1",
		col.KV("url", "tcp://46.28.108.197:8883"),
		col.KV("username", "mqtt1"),
		col.KV("password", "----"))
	p.AddWotServer("conas-dth-esp8266-1", dhtModel, "/conas/dth-esp8266-1", "SIMPLE_URL_ENCODER", "mqtt-1", []string{"http-1"})
	p.Start().Wait()
}
