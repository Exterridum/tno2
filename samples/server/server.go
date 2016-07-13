package main

import (
	"github.com/conas/tno2/wot/model"
	"github.com/conas/tno2/wot/protocol"
)

func main() {
	http := protocol.Http(8080)
	http.Bind("/sample-thing", model.Load("file://../../wot/model/testdata/case-1.json"))
	http.Start()
}
