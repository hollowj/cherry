package timer

const TicketEventName = "TicketEvent"

type TicketEvent struct {
	Id       int64
	CallBack func(ticker *Ticker)
	*Ticker
}

func (*TicketEvent) Name() string {
	return TicketEventName
}

func (p *TicketEvent) UniqueId() int64 {
	return p.Id
}

func NewTicketEvent(actorId int64, callback func(*Ticker), t *Ticker) *TicketEvent {
	event := &TicketEvent{
		Id:       actorId,
		CallBack: callback,
		Ticker:   t,
	}
	return event
}
