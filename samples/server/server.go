package main

import (
	"github.com/conas/tno2/samples/driver"
	"github.com/conas/tno2/wot"
	"github.com/conas/tno2/wot/protocol"
)

func main() {
	server := wot.CreateFromDescriptionUri("file://../../wot/model/testdata/case-1.json")
	driver := driver.ReplyDriver()
	server.Bind(driver)

	http := protocol.Http(8080)
	http.Bind("/sample-thing", server)
	http.Start()
}
