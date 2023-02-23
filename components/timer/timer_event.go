package timer

const TimerEventName = "TimerEvent"

type TimerEvent struct {
	Id       int64
	CallBack func(*Timer)
	*Timer
}

func (*TimerEvent) Name() string {
	return TimerEventName
}

func (p *TimerEvent) UniqueId() int64 {
	return p.Id
}

func NewTimerEvent(actorId int64, callback func(*Timer), t *Timer) *TimerEvent {
	event := &TimerEvent{
		Id:       actorId,
		CallBack: callback,
		Timer:    t,
	}
	return event
}
