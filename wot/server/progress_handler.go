package server

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/sec"
)

type TaskStatusCode int

const (
	TASK_FAILED    TaskStatusCode = -1
	TASK_SCHEDULED TaskStatusCode = 0
	TASK_RUNNING   TaskStatusCode = 1
	TASK_DONE      TaskStatusCode = 2
)

type TaskStatus struct {
	Name      string         `json:"name,omitempty"`
	Status    TaskStatusCode `json:"status"`
	Timestamp time.Time      `json:"timestamp,omitempty"`
	Data      interface{}    `json:"data"`
}

// WotProgressHandler implements async.ProgressHandler
type WotProgressHandler struct {
	name        string
	state       *atomic.Value
	subscribers *async.FanOut
}

func NewWotProgressHandler(name string, state *atomic.Value, subscribers *async.FanOut) *WotProgressHandler {
	return &WotProgressHandler{
		name:        name,
		state:       state,
		subscribers: subscribers,
	}
}

func (ph *WotProgressHandler) Schedule(data interface{}) {
	status := &TaskStatus{
		Name:      ph.name,
		Status:    TASK_SCHEDULED,
		Timestamp: time.Now(),
		Data:      data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *WotProgressHandler) Update(data interface{}) {
	status := &TaskStatus{
		Name:      ph.name,
		Status:    TASK_RUNNING,
		Timestamp: time.Now(),
		Data:      data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *WotProgressHandler) Done(data interface{}) {
	status := &TaskStatus{
		Name:      ph.name,
		Status:    TASK_DONE,
		Timestamp: time.Now(),
		Data:      data,
	}

	ph.state.Store(status)
	ph.subscribers.Publish(status)
}

func (ph *WotProgressHandler) Fail(data interface{}) {
	status := &TaskStatus{
		Name:      ph.name,
		Status:    TASK_FAILED,
		Timestamp: time.Now(),
		Data:      data,
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
