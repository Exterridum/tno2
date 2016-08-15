package main

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/server"
)

func main() {
	// Create WotServer
	wotServer := server.CreateFromDescriptionUri("file://../example-model.json")
	// Attach functionality to WotServer
	setupWotServer(wotServer)

	// WotServer is decoupled from frontend transport protocols. In this step new transport server is created.
	// We can attach any number of generic WotServers under one transport server
	feCfg := col.AsMap(col.KV("port", 8080))
	fe := frontend.NewHTTP(feCfg)
	fe.Bind("/example1", wotServer)
	fe.Start()
}

var db = make(map[string]interface{})

type Throtle struct {
	ThrotlePosition int `json:"throtle-position"`
}

type CriticalEvent struct {
	EventData string `json:"eventData"`
}

// Following section describes WotServer behaviour upon receiving specific requests
func setupWotServer(s *server.WotServer) {
	// Each device can have properties which alter state of the device.
	s.OnGetProperty("relay", func() interface{} {
		// function to access real device/thing properties, such as temperature, etc.
		log.Info("OnGetProperty: relay")
		return db["relay"]
	}).OnUpdateProperty("relay", func(newValue interface{}) {
		// Define how we load property from real device/thing
		log.Info("OnUpdateProperty: relay")
		db["relay"] = newValue
	}).OnInvokeAction("throtle-move", func(args interface{}, ph async.ProgressHandler) interface{} {
		// Define how to interract with actuator. In this case we open throtle of some device/thing
		// Action receives StatusHandler structure which can be used to monitor action progress.
		// Using StatusHandler is not mandatory, but it is a good way to notify users about state of the actions

		// Validate input
		log.Info("OnInvokeAction: throtle-move, position: ", args)

		m := args.(map[string]interface{})

		targetPos := int(m["value"].(float64))
		if targetPos < 0 || targetPos > 50 {
			ph.Fail("Invalid throtle position.")
			return nil
		}

		steps := 4
		step := targetPos / steps

		// Slowly open the throtle
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
