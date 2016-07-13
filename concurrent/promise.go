package concurent

type Promise struct {
	pch chan interface{}
}

func (prev *Promise) Then(stage func(response interface{}) interface{}) *Promise {
	next := Promise{make(chan interface{})}

	go func() {
		next.pch <- stage(<-prev.pch)
		close(next.pch)
	}()

	return &next
}

func NewPromise(c chan interface{}) *Promise {
	return &Promise{c}
}

func Calculate(stage func() interface{}) *Promise {
	p := NewPromise(make(chan interface{}))

	go func() {
		p.pch <- stage()
	}()

	return p
}
