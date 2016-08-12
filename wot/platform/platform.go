package platform

import "github.com/conas/tno2/util/col"

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

func NewComponentConfig(componentType string, cfgParams ...*col.KeyValue) *ComponentConfig {
	params := make(map[string]interface{})
	for _, cfg := range cfgParams {
		params[cfg.K] = cfg.V
	}

	return &ComponentConfig{
		componentType: componentType,
		params:        params,
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
