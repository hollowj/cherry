package timer

const CronEventName = "CronEvent"

type CronEvent struct {
	Id       int64
	CallBack func(cron *Cron)
	*Cron
}

func (*CronEvent) Name() string {
	return CronEventName
}

func (p *CronEvent) UniqueId() int64 {
	return p.Id
}

func NewCronEvent(actorId int64, callback func(cron *Cron), t *Cron) *CronEvent {
	event := &CronEvent{
		Id:       actorId,
		CallBack: callback,
		Cron:     t,
	}
	return event
}
