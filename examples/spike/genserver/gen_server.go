package main

import (
	"log"

	"github.com/conas/tno2/util/async"
)

const (
	MSG_1 async.MessageType = iota
	MSG_2
)

func main() {
	gs := async.NewGenServer().
		HandleCall(MSG_1, func(arg interface{}) interface{} {
			return arg
		}).
		HandleCall(MSG_2, func(arg interface{}) interface{} {
			panic("MSG_2 -> panic")
		})

	gs.Start()

	log.Printf("Output -> %v", gs.Call(MSG_1, 1).Get())
	log.Printf("Output -> %v", gs.Call(MSG_2, 2).Get())
	log.Printf("Output -> %v", gs.Call(MSG_2, 3).Get())
	log.Printf("Output -> %v", gs.Call(MSG_2, 4).Get())
	log.Printf("Output -> %v", gs.Call(MSG_2, 5).Get())
	log.Printf("Output -> %v", gs.Call(MSG_2, 6).Get())
	log.Printf("Output -> %v", gs.Call(MSG_1, 7).Get())
}
