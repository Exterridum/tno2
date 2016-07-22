package main

import (
	"log"
	"time"
)

type Actor struct {
	processor func(<-chan interface{})
	io        chan interface{}
	onPanic   func(Actor, interface{})
}

func Spawn(processor func(<-chan interface{})) chan<- interface{} {
	actor := Actor{
		processor: processor,
		io:        make(chan interface{}),
		onPanic:   restart,
	}

	go actor.read()

	return actor.io
}

func (a Actor) read() {
	defer a.panicHandler()
	a.processor(a.io)
}

func restart(a Actor, message interface{}) {
	a.read()
}

func (a Actor) panicHandler() {
	if err := recover(); err != nil {
		a.onPanic(a, err)
	}
}

func processor(in <-chan interface{}) {
	for {
		mail := <-in
		log.Printf("Agent: message received: %v", mail)

		message := mail.(string)

		if message == "fail" {
			panic("Agent ordered to fail.")
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	a1 := Spawn(processor)

	for {
		a1 <- ("msg1")
		a1 <- ("msg2")
		a1 <- ("fail")
		log.Printf("Loop restart")
	}
}
