package concurent

type Promise struct {
	pch chan interface{}
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

func NewPromise() *Promise {
	return &Promise{
		pch: make(chan interface{}),
	}
}

func Async(callback func() interface{}) *Promise {
	p := NewPromise()

	go func() {
		p.pch <- callback()
	}()

	return p
}
