package main

import (
	"time"

	"github.com/conas/tno2/wot/server"
)

type SampleDevice struct {
	db map[string]interface{}
}

type Throtle struct {
	ThrotlePosition int `json:"throtle-position"`
}

func NewSampleDevice() *SampleDevice {
	return &SampleDevice{
		db: make(map[string]interface{}),
	}
}

func (d *SampleDevice) Init(initParams map[string]interface{}, s *server.WotServer) {

	s.OnGetProperty("relay", func() interface{} {
		return d.db["relay"]
	}).OnUpdateProperty("relay", func(newValue interface{}) {
		d.db["relay"] = newValue
	}).OnInvokeAction("throtle-open", func(position interface{}, status concurent.StatusHandler) {
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

}
