package concurent

type Promise struct {
	pch chan interface{}
}

func (prev *Promise) Then(stage func(response interface{}) interface{}) *Promise {
	next := Promise{make(chan interface{})}

	go func() {
		defer close(prev.pch)
		next.pch <- stage(<-prev.pch)
	}()

	return &next
}

func NewPromise(c chan interface{}) *Promise {
	return &Promise{c}
}

func (p *Promise) End() {
	defer close(p.pch)
	<-p.pch
}

func Calculate(stage func() interface{}) *Promise {
	p := NewPromise(make(chan interface{}))

	go func() {
		p.pch <- stage()
	}()

	return p
}
