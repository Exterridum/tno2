package main

import (
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/server"
)

func main() {
	// Create WotServer and Device
	wotServer := server.CreateFromDescriptionUri("")
	setupServer(wotServer)

	// Attach Device to WotServer
	// This step generic WotServer calls Init method of Device. At the end of this step
	// Wotserver is populated with callbacks to handle Wotserver API calls

	// WotServer is decopled from transport protocols. In this step new transport server is created.
	// We can attach any number of generic WotServers under one transport server
	server.HttpFrontend(8080).Bind("/example1", wotServer).Start()
}

// SampleDevice (Driver) encapsulates communication logic with physical device/thing. As such Device acts as
// translation layer between generit WoTServer and physical device/thing
var db = make(map[string]interface{})

type Throtle struct {
	ThrotlePosition int `json:"throtle-position"`
}

type CriticalEvent struct {
	EventData string `json:"eventData"`
}

func setupServer(s *server.WotServer) {

	// Following section describes WotServer behaviour upon receiving specific requests
	// Each device can have properties which alter state of the device.
	s.OnGetProperty("relay", func() interface{} {
		// Code function to access real device/thing properties, such as temperature, etc.
		return db["relay"]
	}).OnUpdateProperty("relay", func(newValue interface{}) {
		// Define how we load proeprty from real device/thing
		db["relay"] = newValue
	}).OnInvokeAction("throtle-open", func(position interface{}, ph async.ProgressHandler) interface{} {
		// Programm how to interract with actuator. In this case we lineary open throtle of some device/thing
		// Action receives StatusHandler structure which can be used to monitor action progress.
		// Using StatusHandler is not mandatory, but it is a good way to notofy users about state of teh actions

		// Validate input
		targetPos := position.(int)
		if targetPos < 0 || targetPos > 50 {
			ph.Fail("Invalid throtle position.")
			return nil
		}

		steps := 10
		step := position.(int) / steps

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
