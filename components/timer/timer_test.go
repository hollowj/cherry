package timer

import (
	"fmt"
	"testing"
	"time"

	"github.com/cherry-game/cherry/components/timer/dep"
	cherryLogger "github.com/cherry-game/cherry/logger"
)

func TestTimer(t *testing.T) {
	var uid int64 = 1
	AfterFunc(uid, time.Second*5, func(timer *Timer) {
		cherryLogger.Debug("5s after")
	})
	TickerFunc(uid, time.Second*2, func(timer *Ticker) {
		cherryLogger.Debug("2s TickerFunc")
	})
	expr, err := dep.NewCronExpr("* * * * * *")
	if err != nil {
		cherryLogger.Error(err)
		return
	}
	i := 0
	CronFunc(uid, expr, func(timer *Cron) {
		i++
		cherryLogger.Debug("CronFunc")
		if i == 3 {
			timer.Cancel()
		}
	})
}
func TestTimer1(t *testing.T) {
	ttimer := NewEasyTimer()

	var taskChan = make(chan *TimerTask)
	go func() {
		select {
		case c := <-taskChan:
			c.timerChan <- c.callback()
		}
	}()
	task := NewTimerTask(func() ITimer {
		return ttimer.AfterFunc(time.Second, func(timer *Timer) {
			fmt.Println(111)
		})
	})
	taskChan <- task
	cc := <-task.timerChan
	fmt.Println(cc)
}
