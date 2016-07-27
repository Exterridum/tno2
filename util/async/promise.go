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

type ProgressPromise struct {
	pch chan interface{}
	ph  ProgressHandler
}

type ProgressHandler interface {
	Schedule(interface{})
	Update(interface{})
	Done(interface{})
	Fail(interface{})
}

// type StatusHandler func(TaskStatus, interface{})

func RunProgress(task func(ProgressHandler) interface{}, ph ProgressHandler) *ProgressPromise {
	p := NewProgressPromise(ph)

	go func() {
		p.pch <- task(p.ph)
	}()

	return p
}

func NewProgressPromise(ph ProgressHandler) *ProgressPromise {
	return &ProgressPromise{
		pch: make(chan interface{}),
		ph:  ph,
	}
}

func (prev *ProgressPromise) Then(callback func(interface{}, ProgressHandler) interface{}) *ProgressPromise {
	next := NewProgressPromise(prev.ph)

	go func() {
		next.pch <- callback(<-prev.pch, prev.ph)
	}()

	return next
}

func (prev *ProgressPromise) Get() interface{} {
	return <-prev.pch
}
