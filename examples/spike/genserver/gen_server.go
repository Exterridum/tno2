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

	log.Printf("Output -> %v", <-gs.Call(MSG_1, 1))
	log.Printf("Output -> %v", <-gs.Call(MSG_2, 2))
	log.Printf("Output -> %v", <-gs.Call(MSG_2, 3))
	log.Printf("Output -> %v", <-gs.Call(MSG_2, 4))
	log.Printf("Output -> %v", <-gs.Call(MSG_2, 5))
	log.Printf("Output -> %v", <-gs.Call(MSG_2, 6))
	log.Printf("Output -> %v", <-gs.Call(MSG_1, 7))

	time.Sleep(time.Hour * 12)
}
