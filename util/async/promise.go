package async

// ----- Simple Promise

type Promise struct {
	pch chan interface{}
}

func Run(task func() interface{}) *Promise {
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
	statusHandler StatusHandler
}

type TaskStatusCode int

const (
	TASK_SCHEDULED TaskStatusCode = iota
	TASK_RUNNING
	TASK_DONE
	TASK_FAILED
)

type TaskStatus struct {
	Code TaskStatusCode `json:"code"`
	Data interface{}    `json:"data"`
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
		p.Channel() <- task(p.statusHandler)
	}()

	return p
}

func NewStatusPromise(statusHandler StatusHandler) *StatusPromise {
	return &StatusPromise{
		pch:           make(chan interface{}),
		statusHandler: statusHandler,
	}
}

func (p *StatusPromise) Channel() chan<- interface{} {
	return p.pch
}

func (prev *StatusPromise) Then(callback func(interface{}, StatusHandler) interface{}) *StatusPromise {
	next := NewStatusPromise(prev.statusHandler)

	go func() {
		next.Channel() <- callback(<-prev.pch, prev.statusHandler)
	}()

	return next
}

func (prev *StatusPromise) Wait() interface{} {
	return <-prev.pch
}
