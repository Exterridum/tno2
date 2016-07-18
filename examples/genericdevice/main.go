package main

import (
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/server/protocol"
)

var modelPrefix = "file://models/"
var models = []string{"case-property-001", "case-action-001"}

func main() {
	http := protocol.NewHttp(8080)
	bindSyncModels(http)
	http.Start()
}

func bindSyncModels(http *protocol.Http) {
	for _, model := range models {
		server := wot.CreateFromDescriptionUri(modelURI(model))
		device := NewTestDevice()
		server.Connect(device, make(map[string]interface{}))
		http.Bind(str.Concat("/", model), server)
	}
}

func modelURI(model string) string {
	return str.Concat(modelPrefix, model, ".json")
}
