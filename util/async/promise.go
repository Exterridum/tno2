package async

// ----- Simple Promise

type Promise struct {
	pch chan interface{}
}

func Run(task func() interface{}) *Promise {
	p := NewPromise()

	go func() {
		p.pch <- task()
	}()

	return p
}

func NewPromise() *Promise {
	return &Promise{
		pch: make(chan interface{}),
	}
}

func (prev *Promise) Then(callback func(response interface{}) interface{}) *Promise {
	next := NewPromise()

	go func() {
		next.pch <- callback(<-prev.pch)
	}()

	return next
}

func (prev *Promise) Get() interface{} {
	return <-prev.pch
}

// ----- Promise With Status Update

type StatusPromise struct {
	pch           chan interface{}
	statusHandler StatusHandler
}

type StatusHandler interface {
	Schedule(interface{})
	Update(interface{})
	Done(interface{})
	Fail(interface{})
}

// type StatusHandler func(TaskStatus, interface{})

func RunWithStatus(task func(StatusHandler) interface{}, statusHandler StatusHandler) *StatusPromise {
	p := NewStatusPromise(statusHandler)

	go func() {
		p.pch <- task(p.statusHandler)
	}()

	return p
}

func NewStatusPromise(statusHandler StatusHandler) *StatusPromise {
	return &StatusPromise{
		pch:           make(chan interface{}),
		statusHandler: statusHandler,
	}
}

func (prev *StatusPromise) Then(callback func(interface{}, StatusHandler) interface{}) *StatusPromise {
	next := NewStatusPromise(prev.statusHandler)

	go func() {
		next.pch <- callback(<-prev.pch, prev.statusHandler)
	}()

	return next
}

func (prev *StatusPromise) Get() interface{} {
	return <-prev.pch
}
