package main

import (
	"log"
	"time"
)

func processor(in <-chan interface{}) {
	log.Printf("Entering processor")
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
	go func() {
		for {
		}
	}()

	c1 := make(chan interface{})

	c1 <- 1

	log.Println("writen")

	for {
	}

	// a1 := async.Spawn(processor).Channel()

	// for {
	// 	a1 <- ("msg1")
	// 	a1 <- ("msg2")
	// 	a1 <- ("fail")
	// }
}
