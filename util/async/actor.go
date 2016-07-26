package async

type Actor struct {
	processor func(<-chan interface{})
	io        chan interface{}
	onPanic   func(interface{})
}

func Spawn(processor func(<-chan interface{}), panicHandler func(interface{})) *Actor {
	actor := &Actor{
		processor: processor,
		io:        make(chan interface{}),
		onPanic:   panicHandler,
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
			if a.onPanic != nil {
				a.onPanic(err)
			}

			a.read()
		}
	}()

	a.processor(a.io)
}
