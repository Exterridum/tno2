package server

import (
	"sync/atomic"

	"github.com/conas/tno2/util/async"
)

type ActionStatusHandler struct {
	Value *atomic.Value
}

func (ash *ActionStatusHandler) Schedule(data interface{}) {
	ash.Value.Store(&async.TaskStatus{
		Code: async.TASK_SCHEDULED,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Update(data interface{}) {
	ash.Value.Store(&async.TaskStatus{
		Code: async.TASK_RUNNING,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Done(data interface{}) {
	ash.Value.Store(&async.TaskStatus{
		Code: async.TASK_DONE,
		Data: data,
	})
}

func (ash *ActionStatusHandler) Fail(data interface{}) {
	ash.Value.Store(&async.TaskStatus{
		Code: async.TASK_FAILED,
		Data: data,
	})
}

func newActionStatusHandler() ActionStatusHandler {
	return ActionStatusHandler{
		Value: new(atomic.Value),
	}
}
