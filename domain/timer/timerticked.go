package timer

import (
	"encoding/json"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	TimerTicked interface {
		event.Event
	}

	timerTicked struct{}
)

func NewTimerTicked() TimerTicked {
	return &timerTicked{}
}

func (e *timerTicked) Type() event.EventType {
	return event.EventTypeTimerTicked
}

func (e *timerTicked) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{
		e.Type().String(),
	})
}

func (e *timerTicked) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
