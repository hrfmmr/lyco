package timer

import "github.com/hrfmmr/lyco/domain/event"

type (
	TimerFinished interface {
		event.Event
	}

	timerFinished struct{}
)

func NewTimerFinished() TimerFinished {
	return &timerFinished{}
}

func (e *timerFinished) Type() event.EventType {
	return event.EventTypeTimerFinished
}
