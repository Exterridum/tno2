package main

import (
	"time"

	"github.com/conas/tno2/wot/backend"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/server"
)

var refModel = "file://reference-model.json"

func main() {
	p := NewPlatform()

	fe1 := NewComponentConfig("HTTP", KV("port", 8080))
	p.AddFronted("http-1", fe1)

	be1 := NewComponentConfig("MQTT", KV("url", "tcp://localhost"), KV("port", 1883))
	p.AddBackend("mqtt-1", be1)

	simulator := NewWotServer("file://reference-model.json", "mqtt-1", "http-1")
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

type Platform struct {
	frontends map[string]*ComponentConfig
	backends  map[string]*ComponentConfig
	wots      map[string]*WotConfig
}

func NewPlatform() *Platform {
	return &Platform{
		frontends: make(map[string]*ComponentConfig),
		backends:  make(map[string]*ComponentConfig),
		wots:      make(map[string]*WotConfig),
	}
}

func (p *Platform) AddFronted(id string, cc *ComponentConfig) {
	p.frontends[id] = cc
}

func (p *Platform) AddBackend(id string, cc *ComponentConfig) {
	p.backends[id] = cc
}

func (p *Platform) AddWotServer(id string, wotConfig *WotConfig) {
	p.wots[id] = wotConfig
}

func (p *Platform) Start() {
}

type ComponentConfig struct {
	componentType string
	params        map[string]interface{}
}

func NewComponentConfig(componentType string, cfgParams ...*KeyValue) *ComponentConfig {
	params := make(map[string]interface{})
	for _, cfg := range cfgParams {
		params[cfg.k] = cfg.v
	}

	return &ComponentConfig{
		componentType: componentType,
		params:        params,
	}
}

type KeyValue struct {
	k string
	v interface{}
}

func KV(k string, v interface{}) *KeyValue {
	return &KeyValue{
		k: k,
		v: v,
	}
}

type WotConfig struct {
	modelUri  string
	frontends []string
	backends  []string
}

func NewWotServer(uri, backend string, frontends ...string) *WotConfig {
	return nil
}
