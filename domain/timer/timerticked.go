package timer

import (
	"encoding/json"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	TimerTicked interface {
		event.Event
		Mode() TimerMode
	}

	timerTicked struct {
		mode TimerMode
	}
)

func NewTimerTicked(mode TimerMode) TimerTicked {
	return &timerTicked{mode}
}

func (e *timerTicked) Type() event.EventType {
	return event.EventTypeTimerTicked
}

func (e *timerTicked) Mode() TimerMode {
	return e.mode
}

func (e *timerTicked) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
		Mode string `json:"mode"`
	}{
		e.Type().String(),
		e.Mode().String(),
	})
}

func (e *timerTicked) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
