package server

import (
	"sync"
	"sync/atomic"

	"github.com/conas/tno2/util/async"
	"github.com/conas/tno2/util/sec"
)

type ActionStatusHandler struct {
	slot *atomic.Value
}

func NewActionStatusHandler(slot *atomic.Value) *ActionStatusHandler {
	return &ActionStatusHandler{
		slot: slot,
	}
}

func (ash *ActionStatusHandler) Schedule(data interface{}) {
	ash.slot.Store(&async.TaskStatus{
		Code: async.TASK_SCHEDULED,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Update(data interface{}) {
	ash.slot.Store(&async.TaskStatus{
		Code: async.TASK_RUNNING,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Done(data interface{}) {
	ash.slot.Store(&async.TaskStatus{
		Code: async.TASK_DONE,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Fail(data interface{}) {
	ash.slot.Store(&async.TaskStatus{
		Code: async.TASK_FAILED,
		Data: data,
	})
}

type ActionResults struct {
	rwmut *sync.RWMutex
	slots map[string]*atomic.Value
}

func NewActionResults() *ActionResults {
	return &ActionResults{
		rwmut: &sync.RWMutex{},
		slots: make(map[string]*atomic.Value),
	}
}

func (ar *ActionResults) CreateSlot() (string, *atomic.Value) {
	slotID, _ := sec.UUID4()

	ar.rwmut.Lock()
	defer ar.rwmut.Unlock()

	ar.slots[slotID] = &atomic.Value{}

	return slotID, ar.slots[slotID]
}

func (ar *ActionResults) GetSlot(slotID string) (*atomic.Value, bool) {
	ar.rwmut.RLock()
	defer ar.rwmut.RUnlock()

	slot, rc := ar.slots[slotID]

	return slot, rc
}
