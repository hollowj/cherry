package timer

import (
	"runtime/debug"
	"time"

	"github.com/cherry-game/cherry"
	"github.com/cherry-game/cherry/components/timer/dep"
	cherryLogger "github.com/cherry-game/cherry/logger"
)

type IEasyTimer interface {
	AfterFunc(d time.Duration, cb func(*Timer)) *Timer
	CronFunc(cronExpr *dep.CronExpr, cb func(cron *Cron)) *Cron
	TickerFunc(d time.Duration, cb func(ticker *Ticker)) *Ticker
	Destroy()
}
type EasyTimer struct {
	dispatcher *Dispatcher //timer

	mapActiveTimer map[ITimer]struct{}
}

func NewEasyTimer() *EasyTimer {
	return &EasyTimer{
		mapActiveTimer: make(map[ITimer]struct{}),
		dispatcher:     NewDispatcher(1000),
	}
}

func (s *EasyTimer) GetDispatcher() *Dispatcher {
	return s.dispatcher
}
func (s *EasyTimer) OnCloseTimer(t ITimer) {
	delete(s.mapActiveTimer, t)
}

func (s *EasyTimer) OnAddTimer(t ITimer) {
	if t != nil {
		if s.mapActiveTimer == nil {
			s.mapActiveTimer = map[ITimer]struct{}{}
		}

		s.mapActiveTimer[t] = struct{}{}
	}
}
func (s *EasyTimer) AfterFunc(d time.Duration, cb func(*Timer)) *Timer {
	return s.dispatcher.AfterFunc(d, nil, cb, s.OnCloseTimer, s.OnAddTimer)
}
func (s *EasyTimer) CronFunc(cronExpr *dep.CronExpr, cb func(cron *Cron)) *Cron {
	return s.dispatcher.CronFunc(cronExpr, nil, cb, s.OnCloseTimer, s.OnAddTimer)

}
func (s *EasyTimer) TickerFunc(d time.Duration, cb func(ticker *Ticker)) *Ticker {
	safeFun := func(ticker *Ticker) {
		defer func() {
			if err := recover(); err != nil {
				cherryLogger.Errorf("%v", err)
				stack := string(debug.Stack())
				cherryLogger.Error(stack)
			}
		}()
		cb(ticker)
	}
	return s.dispatcher.TickerFunc(d, nil, safeFun, s.OnCloseTimer, s.OnAddTimer)
}

func (s *EasyTimer) Destroy() {
	for pTimer := range s.mapActiveTimer {
		pTimer.Cancel()
	}

	s.mapActiveTimer = nil
}

func AfterFunc(uid int64, d time.Duration, cb func(timer *Timer)) ITimer {
	task := NewTimerTask(func() ITimer {
		return _timer.AfterFunc(d, func(t *Timer) {
			cherry.PostEvent(NewTimerEvent(uid, cb, t))
		})
	})
	taskChan <- task
	t := <-task.timerChan
	return t
}
func CronFunc(uid int64, cronExpr *dep.CronExpr, cb func(cron *Cron)) ITimer {
	task := NewTimerTask(func() ITimer {
		return _timer.CronFunc(cronExpr, func(t *Cron) {
			cherry.PostEvent(NewCronEvent(uid, cb, t))
		})
	})
	taskChan <- task
	t := <-task.timerChan
	return t
}
func TickerFunc(uid int64, d time.Duration, cb func(ticker *Ticker)) ITimer {
	task := NewTimerTask(func() ITimer {
		return _timer.TickerFunc(d, func(t *Ticker) {
			cherry.PostEvent(NewTicketEvent(uid, cb, t))
		})
	})
	taskChan <- task
	t := <-task.timerChan
	return t
}
