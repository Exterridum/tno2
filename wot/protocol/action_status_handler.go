package protocol

import (
	"sync/atomic"

	"github.com/conas/tno2/util/concurrent"
)

type ActionStatusHandler struct {
	Value *atomic.Value
}

func (ash *ActionStatusHandler) Schedule(data interface{}) {
	ash.Value.Store(&concurent.TaskStatus{
		Code: concurent.SCHEDULED,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Update(data interface{}) {
	ash.Value.Store(&concurent.TaskStatus{
		Code: concurent.RUNNING,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Done(data interface{}) {
	ash.Value.Store(&concurent.TaskStatus{
		Code: concurent.DONE,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Fail(data interface{}) {
	ash.Value.Store(&concurent.TaskStatus{
		Code: concurent.FAILED,
		Data: data,
	})
}

func newActionStatusHandler() ActionStatusHandler {
	return ActionStatusHandler{
		Value: new(atomic.Value),
	}
}
