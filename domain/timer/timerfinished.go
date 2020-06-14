package timer

import (
	"encoding/json"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	TimerFinished interface {
		event.Event
		Mode() TimerMode
	}

	timerFinished struct {
		mode TimerMode
	}
)

func NewTimerFinished(mode TimerMode) TimerFinished {
	return &timerFinished{mode}
}

func (e *timerFinished) Type() event.EventType {
	return event.EventTypeTimerFinished
}

func (e *timerFinished) Mode() TimerMode {
	return e.mode
}

func (e *timerFinished) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Mode string `json:"mode"`
	}{
		e.Type().String(),
		e.Mode().String(),
	})
}

func (e *timerFinished) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
