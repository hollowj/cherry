package timer

import (
	"time"

	cfacade "github.com/cherry-game/cherry/facade"
)

const (
	Component_TIMER = "timer"
)

type (
	TimerComponent struct {
		cfacade.Component
	}
)

var _timer *EasyTimer
var taskChan chan *TimerTask

func NewTimerComponent() *TimerComponent {
	return &TimerComponent{}
}
func (r *TimerComponent) Init() {
	taskChan = make(chan *TimerTask, 1000)
	StartTimer(time.Millisecond*50, 100000)

	_timer = NewEasyTimer()
	go func() {
		for {
			select {
			case t := <-_timer.GetDispatcher().ChanTimer:
				t.Do()
			case task := <-taskChan:
				task.timerChan <- task.callback()
			}
		}
	}()
}
func (d *TimerComponent) Name() string {
	return Component_TIMER
}
func (d *TimerComponent) OnStop() {
	_timer.Destroy()
}
