package cherryHandler

import (
	"runtime/debug"

	cfacade "github.com/cherry-game/cherry/facade"
	clog "github.com/cherry-game/cherry/logger"
)

type (
	ExecutorEvent struct {
		Executor
		eventData  cfacade.IEvent
		eventSlice []cfacade.EventFn
	}
)

func (p *ExecutorEvent) EventData() cfacade.IEvent {
	return p.eventData
}

func (p *ExecutorEvent) Invoke() {
	defer func() {
		if rev := recover(); rev != nil {
			clog.Warnf("recover in Event.%v %s", rev, string(debug.Stack()))
			clog.Warnf("event = [%+v]", p.eventData)
		}
	}()

	for _, fn := range p.eventSlice {
		fn(p.eventData)
	}
}

func (p *ExecutorEvent) QueueHash(queueNum int) int {
	return int(p.eventData.UniqueId() % int64(queueNum))
}
