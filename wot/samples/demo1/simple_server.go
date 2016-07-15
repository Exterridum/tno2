package main

import (
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/protocol"
)

func main() {
	server := wot.CreateFromDescriptionUri("file://../../model/testdata/case-1.json")
	driver := NewDemoDriver()
	server.BindSync(driver)

	http := protocol.NewHttp(8080)
	http.Bind("/sample-thing", server)
	http.Start()
}
