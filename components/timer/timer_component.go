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

func NewTimerComponent() *TimerComponent {
	return &TimerComponent{}
}
func (r *TimerComponent) Init() {
	StartTimer(time.Millisecond*50, 1000000)

}
func (d *TimerComponent) Name() string {
	return Component_TIMER
}
