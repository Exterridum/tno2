package main

import (
	"time"

	"github.com/conas/tno2/util/async"
)

func processor(in <-chan interface{}) {
	for {
		mail := <-in
		message := mail.(bool)

		if message == true {
			panic("Actor fail.")
		}

		time.Sleep(2 * time.Second)
	}
}

func main() {
	a1 := async.Spawn(processor, nil).Channel()

	for {
		a1 <- (true)
	}
}
