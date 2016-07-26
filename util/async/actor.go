package async

import "log"

type Actor struct {
	processor func(<-chan interface{})
	io        chan interface{}
}

func Spawn(processor func(<-chan interface{})) *Actor {
	actor := &Actor{
		processor: processor,
		io:        make(chan interface{}),
	}

	go actor.read()

	return actor
}

func (a *Actor) Channel() chan<- interface{} {
	return a.io
}

func (a *Actor) read() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Actor panick -> %v", err)
			a.read()
		}
	}()

	a.processor(a.io)
}
