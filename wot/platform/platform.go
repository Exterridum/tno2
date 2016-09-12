package platform

import (
	"sync"

	"github.com/conas/tno2/util/col"
	"github.com/conas/tno2/wot/backend"
	"github.com/conas/tno2/wot/frontend"
	"github.com/conas/tno2/wot/server"
)

var feTypes map[string]frontend.Factory = make(map[string]frontend.Factory)
var beTypes map[string]backend.Factory = make(map[string]backend.Factory)

type Platform struct {
	hostname  string
	frontends map[string]frontend.Frontend
	backends  map[string]backend.Backend
	wots      map[string]*server.WotServer
}

func init() {
	RegisterFrontendType("HTTP", frontend.NewHTTP)
	RegisterBackendType("MQTT-1", backend.NewMQTT_1)
}

func NewPlatform(hostname string) *Platform {
	return &Platform{
		hostname:  hostname,
		frontends: make(map[string]frontend.Frontend),
		backends:  make(map[string]backend.Backend),
		wots:      make(map[string]*server.WotServer),
	}
}

func RegisterFrontendType(feTypeID string, factory frontend.Factory) {
	feTypes[feTypeID] = factory
}

func RegisterBackendType(beTypeID string, factory backend.Factory) {
	beTypes[beTypeID] = factory
}

func (p *Platform) AddFrontend(feID, feType string, cfgParams ...*col.KeyValue) {
	params := make(map[string]interface{})
	for _, cfg := range cfgParams {
		params[cfg.K] = cfg.V
	}
	params["hostname"] = p.hostname

	fe := feTypes[feType](params)
	p.frontends[feID] = fe
}

func (p *Platform) AddBackend(bedID, beType string, cfgParams ...*col.KeyValue) {
	params := col.AsMap(cfgParams)

	be := beTypes[beType](params)
	p.backends[bedID] = be
}

func (p *Platform) AddWotServer(id, wotDescURI, ctxPath, beEncID, beID string, feIDs []string) {
	wotServer := server.CreateFromDescriptionUri(wotDescURI)
	p.wots[id] = wotServer
	be, _ := p.backends[beID]
	encoder, error := backend.Encoders.Get(beEncID)

	if error != nil {
		panic(error)
	}

	be.Bind(wotServer, ctxPath, encoder)

	for _, feId := range feIDs {
		frontend, _ := p.frontends[feId]
		frontend.Bind(ctxPath, wotServer)
	}
}

func (p *Platform) WotServer(id string) *server.WotServer {
	return p.wots[id]
}

func (p *Platform) Start() *sync.WaitGroup {
	wg := &sync.WaitGroup{}

	for _, fe := range p.frontends {
		wg.Add(1)
		go func() {
			fe.Start()
		}()
	}

	for _, be := range p.backends {
		wg.Add(1)
		go func() {
			be.Start()
		}()
	}

	return wg
}
