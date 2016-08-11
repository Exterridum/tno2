package async

type Actor struct {
	processor func(<-chan interface{})
	in        chan interface{}
	onPanic   func(interface{})
	run       bool
}

func Spawn(
	mailboxSize int,
	processor func(<-chan interface{}),
	panicHandler func(interface{})) *Actor {

	actor := &Actor{
		processor: processor,
		in:        make(chan interface{}, mailboxSize),
		onPanic:   panicHandler,
		run:       true,
	}

	actor.start()

	return actor
}

func (a *Actor) Channel() chan<- interface{} {
	return a.in
}

func (a *Actor) start() {
	go func() {
		for a.run {
			a.read()
		}
	}()
}

func (a *Actor) read() {
	defer a.recovery()
	//a.processor is message processing loop. In case of processor panic
	//a.run = false is bypassed and defered recovery method takes controll
	//calling goroutine the restarts the reading loop
	//In case of processor is finished normally, e.g. channel is closed
	//a.run flag is set to false, function ends and goroutine will not
	//attempt to restart processing loop.
	a.processor(a.in)
	a.run = false
}

func (a *Actor) recovery() {
	if err := recover(); err != nil {
		if a.onPanic != nil {
			a.onPanic(err)
		}
	}
}
