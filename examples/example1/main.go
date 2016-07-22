package main

import (
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/wot/server"
)

func main() {
	// Create WotServer and Device
	wotServer := server.CreateFromDescriptionUri("")
	device := NewSampleDevice()

	// Attach Device to WotServer
	// This step generic WotServer calls Init method of Device. At the end of this step
	// Wotserver is populated with callbacks to handle Wotserver API calls
	wotServer.Connect(device, make(map[string]interface{}))

	// WotServer is decopled from transport protocols. In this step new transport server is created.
	// We can attach any number of generic WotServers under one transport server
	http := server.NewHttp(8080)
	http.Bind("/sample-device", wotServer)

	// Starting transport server
	http.Start()
}

func NewSampleDevice() *SampleDevice {
	return &SampleDevice{
		db: make(map[string]interface{}),
	}
}

// SampleDevice (Driver) encapsulates communication logic with physical device/thing. As such Device acts as
// translation layer between generit WoTServer and physical device/thing
type SampleDevice struct {
	db map[string]interface{}
}

type Throtle struct {
	ThrotlePosition int `json:"throtle-position"`
}

type CriticalEvent struct {
	EventData string `json:"eventData"`
}

func (d *SampleDevice) Init(initParams map[string]interface{}, s *server.WotServer) {

	// Following section describes WotServer behaviour upon receiving specific requests
	// Each device can have properties which alter state of the device.
	s.OnGetProperty("relay", func() interface{} {
		// Code function to access real device/thing properties, such as temperature, etc.
		return d.db["relay"]
	}).OnUpdateProperty("relay", func(newValue interface{}) {
		// Define how we load proeprty from real device/thing
		d.db["relay"] = newValue
	}).OnInvokeAction("throtle-open", func(position interface{}, status async.StatusHandler) {
		// Programm how to interract with actuator. In this case we lineary open throtle of some device/thing
		// Action receives StatusHandler structure which can be used to monitor action progress.
		// Using StatusHandler is not mandatory, but it is a good way to notofy users about state of teh actions

		// Validate input
		targetPos := position.(int)
		if targetPos < 0 || targetPos > 50 {
			status.Fail("Invalid throtle position.")
			return
		}

		steps := 10
		step := position.(int) / steps

		// Slowly open the throtle
		for i := 0; i < steps; i++ {
			status.Update(&Throtle{ThrotlePosition: i * step})
			time.Sleep(time.Second * 5)
		}

		status.Done(&Throtle{ThrotlePosition: targetPos})
	})

	//Event generator. For demo purposes event geneator generates some sample events, clients  can subscribe to.
	go func() {
		for {
			s.EmitEvent("critical-temperature-event", &CriticalEvent{EventData: "temperature -> 95 °C and rasing."})
			time.Sleep(time.Second * 5)
		}
	}()

}
