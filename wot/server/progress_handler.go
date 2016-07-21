package server

import (
	"sync"
	"sync/atomic"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/sec"
)

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

type ProgressHandler struct {
	state       *atomic.Value
	subscribers *async.FanOut
}

func NewProgressHandler(state *atomic.Value, subscribers *async.FanOut) *ProgressHandler {
	return &ProgressHandler{
		state:       state,
		subscribers: subscribers,
	}
}

func (ph *ProgressHandler) Schedule(data interface{}) {
	status := &TaskStatus{
		Code: TASK_SCHEDULED,
		Data: data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *ProgressHandler) Update(data interface{}) {
	status := &TaskStatus{
		Code: TASK_RUNNING,
		Data: data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *ProgressHandler) Done(data interface{}) {
	status := &TaskStatus{
		Code: TASK_DONE,
		Data: data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *ProgressHandler) Fail(data interface{}) {
	status := &TaskStatus{
		Code: TASK_FAILED,
		Data: data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

type ActionResults struct {
	rwmut  *sync.RWMutex
	states map[string]*atomic.Value
}

func NewActionResults() *ActionResults {
	return &ActionResults{
		rwmut:  &sync.RWMutex{},
		states: make(map[string]*atomic.Value),
	}
}

func (ar *ActionResults) CreateSlot() (string, *atomic.Value) {
	stateID, _ := sec.UUID4()

	ar.rwmut.Lock()
	defer ar.rwmut.Unlock()

	ar.states[stateID] = &atomic.Value{}

	return stateID, ar.states[stateID]
}

func (ar *ActionResults) GetSlot(stateID string) (*atomic.Value, bool) {
	ar.rwmut.RLock()
	defer ar.rwmut.RUnlock()

	state, rc := ar.states[stateID]

	return state, rc
}
