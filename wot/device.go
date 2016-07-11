package main

import (
	"github.com/conas/tno2/wot/model"
	"github.com/conas/tno2/wot/protocol"
)

func main() {
	http := protocol.New()
	http.Attach(model.Load("file://model/testdata/case-1.json"))
	http.Start(":8080")
}
