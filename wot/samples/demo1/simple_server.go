package main

import (
	"github.com/conas/tno2/util/strings"
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/protocol"
)

var modelPrefix = "file://../../model/testdata/"
var models = []string{"case-property-001"}

func main() {
	http := protocol.NewHttp(8080)
	bindModels(http)
	http.Start()
}

func bindModels(http *protocol.Http) {
	for _, model := range models {
		server := wot.CreateFromDescriptionUri(modelURI(model))
		driver := NewDemoDriver()
		server.BindSync(driver, make(map[string]interface{}))
		http.Bind(strings.Concat("/", model), server)
	}
}

func modelURI(model string) string {
	return strings.Concat(modelPrefix, model, ".json")
}
