package main

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/server"
)

var model = "file://../example-model.json"

//Basic low level wotServer setup
func main() {
	//WoT server defines implements interaction with device
	wotServer := server.CreateFromDescriptionUri(model)
	setupWotServer(wotServer)

	//Frontend is transport implementation
	feCfg := col.AsMap([]*col.KeyValue{col.KV("port", 8080)})
	fe := frontend.NewHTTP(feCfg)
	fe.Bind("/01-basic-example", wotServer)
	fe.Start()
}

var db = make(map[string]interface{})

type Throtle struct {
	ThrotlePosition int `json:"throtle-position"`
}

type CriticalEvent struct {
	EventData string `json:"eventData"`
}

//Implementation of interaction with Web Device
func setupWotServer(s *server.WotServer) {
	s.OnGetProperty("relay", func() interface{} {
		log.Info("OnGetProperty: relay")
		return db["relay"]
	}).OnUpdateProperty("relay", func(newValue interface{}) {
		log.Info("OnUpdateProperty: relay")
		db["relay"] = newValue
	}).OnInvokeAction("throtle-move", func(args interface{}, ph async.ProgressHandler) interface{} {
		log.Info("OnInvokeAction: throtle-move, position: ", args)
		m := args.(map[string]interface{})
		targetPos := int(m["value"].(float64))

		if targetPos < 0 || targetPos > 50 {
			ph.Fail("Invalid throtle position.")
			return nil
		}

		steps := 4
		step := targetPos / steps

		for i := 0; i < steps; i++ {
			ph.Update(&Throtle{ThrotlePosition: i * step})
			time.Sleep(time.Second * 5)
		}

		return Throtle{ThrotlePosition: targetPos}
	})

	//Event generator. For demo purposes event geneator generates some sample events, clients  can subscribe to.
	go func() {
		for {
			s.EmitEvent("critical-temperature-event", &CriticalEvent{EventData: "temperature -> 192°C and rasing."})
			time.Sleep(time.Second * 5)
		}
	}()
}
