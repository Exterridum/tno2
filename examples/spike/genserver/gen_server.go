package main

import (
	"log"
	"time"

	"github.com/conas/tno2/util/async"
)

const (
	MSG_1 async.MessageType = iota
	MSG_2
)

func main() {
	gs := async.NewGenServer().
		Handle(MSG_1, func(arg interface{}) interface{} {
			v := arg.(int)
			return v + 1
		}).
		Handle(MSG_2, func(arg interface{}) interface{} {
			panic("MSG_2 -> panic")
		})

	gs.Start()

	defer timeTrack(time.Now(), "channels")

	for i := 0; i < 1000000; i++ {
		<-gs.Call(MSG_1, i)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
