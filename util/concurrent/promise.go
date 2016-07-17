package concurent

// ----- Simple Promise

type Promise struct {
	pch chan interface{}
}

func Async(task func() interface{}) *Promise {
	p := NewPromise()

	go func() {
		p.Channel() <- task()
	}()

	return p
}

func NewPromise() *Promise {
	return &Promise{
		pch: make(chan interface{}),
	}
}

func (p *Promise) Channel() chan<- interface{} {
	return p.pch
}

func (prev *Promise) Then(callback func(response interface{}) interface{}) *Promise {
	next := NewPromise()

	go func() {
		next.Channel() <- callback(<-prev.pch)
	}()

	return next
}

func (prev *Promise) Wait() interface{} {
	return <-prev.pch
}

// ----- Promise With Status Update

type StatusPromise struct {
	pch           chan interface{}
	statusHandler *StatusHandler
}

type StatusHandler func(int, string)

func AsyncStatus(task func(*StatusHandler) interface{}, statusHandler StatusHandler) *StatusPromise {
	p := NewStatusPromise(&statusHandler)

	go func() {
		p.Channel() <- task(p.statusHandler)
	}()

	return p
}

func NewStatusPromise(statusHandler *StatusHandler) *StatusPromise {
	return &StatusPromise{
		pch:           make(chan interface{}),
		statusHandler: statusHandler,
	}
}

func (p *StatusPromise) Channel() chan<- interface{} {
	return p.pch
}

func (prev *StatusPromise) Then(callback func(interface{}, *StatusHandler) interface{}) *StatusPromise {
	next := NewStatusPromise(prev.statusHandler)

	go func() {
		next.Channel() <- callback(<-prev.pch, prev.statusHandler)
	}()

	return next
}

func (prev *StatusPromise) Wait() interface{} {
	return <-prev.pch
}
