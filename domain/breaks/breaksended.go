package breaks

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	BreaksEnded interface {
		event.Event
		StartedAt() StartedAt
		Duration() Duration
		EndedAt() EndedAt
	}

	breaksEnded struct {
		startedAt StartedAt
		duration  Duration
		endedAt   EndedAt
	}
)

func NewBreaksEnded(startedAt StartedAt, duration Duration, endedAt EndedAt) BreaksEnded {
	return &breaksEnded{
		startedAt,
		duration,
		endedAt,
	}
}

func (e *breaksEnded) Type() event.EventType {
	return event.EventTypeBreaksEnded
}

func (e *breaksEnded) StartedAt() StartedAt {
	return e.startedAt
}

func (e *breaksEnded) Duration() Duration {
	return e.duration
}

func (e *breaksEnded) EndedAt() EndedAt {
	return e.endedAt
}

func (e *breaksEnded) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string `json:"type"`
		StartedAt string `json:"started_at"`
		Duration  string `json:"duration"`
		EndedAt   string `json:"ended_at"`
	}{
		e.Type().String(),
		time.Unix(0, e.startedAt.Value()).String(),
		time.Duration(e.duration.Value()).String(),
		time.Unix(0, e.endedAt.Value()).String(),
	})
}

func (e *breaksEnded) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
