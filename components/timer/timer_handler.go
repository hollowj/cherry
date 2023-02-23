package timer

import (
	cherryFacade "github.com/cherry-game/cherry/facade"
	chandler "github.com/cherry-game/cherry/net/handler"
)

type (
	TimerHandler struct {
		chandler.Handler
	}
)

func (h *TimerHandler) Name() string {
	return "timerHandler"
}

func (h *TimerHandler) OnInit() {
	h.AddEvent(TimerEventName, onTimerEvent)
	h.AddEvent(CronEventName, onCronEvent)
	h.AddEvent(TicketEventName, onTicketEvent)

}
func onTimerEvent(e cherryFacade.IEvent) {
	timerEvent := e.(*TimerEvent)
	timerEvent.CallBack(timerEvent.Timer)
}
func onCronEvent(e cherryFacade.IEvent) {
	cronEvent := e.(*CronEvent)
	cronEvent.CallBack(cronEvent.Cron)
}
func onTicketEvent(e cherryFacade.IEvent) {
	ticketEvent := e.(*TicketEvent)
	ticketEvent.CallBack(ticketEvent.Ticker)
}
