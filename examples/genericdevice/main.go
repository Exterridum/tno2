package main

import (
	"github.com/conas/tno2/util/str"
	"github.com/conas/tno2/wot/server"
)

var modelPrefix = "file://models/"
var models = []string{"case-property-001", "case-action-001"}

func main() {
	http := server.NewHttp(8080)
	bindSyncModels(http)
	http.Start()
}

func bindSyncModels(http *server.Http) {
	for _, model := range models {
		server := server.CreateFromDescriptionUri(modelURI(model))
		device := NewTestDevice()
		server.Connect(device, make(map[string]interface{}))
		http.Bind(str.Concat("/", model), server)
	}
}

func modelURI(model string) string {
	return str.Concat(modelPrefix, model, ".json")
}
