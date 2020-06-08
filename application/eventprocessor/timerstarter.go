package eventprocessor

import (
	"fmt"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type TimerStarter struct {
	tmr timer.Timer
}

func NewTimerStarter(tmr timer.Timer) *TimerStarter {
	return &TimerStarter{tmr}
}

func (s *TimerStarter) EventType() event.EventType {
	return event.EventTypeTaskStarted
}

func (s *TimerStarter) HandleEvent(e event.Event) {
	ev, ok := e.(task.TaskStarted)
	if !ok {
		panic(fmt.Sprintf("😕 got unexpected event:%v, expecting: task.TaskStarted", e))
	}
	s.tmr.Start(ev.Duration(), ev.Elapsed())
}
