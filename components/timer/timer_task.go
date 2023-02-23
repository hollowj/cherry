package timer

type TimerTask struct {
	timerChan chan ITimer
	callback  func() ITimer
}

func NewTimerTask(f func() ITimer) *TimerTask {
	task := &TimerTask{
		timerChan: make(chan ITimer),
		callback:  f,
	}
	return task
}
