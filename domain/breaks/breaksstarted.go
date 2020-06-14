package breaks

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	BreaksStarted interface {
		event.Event
		StartedAt() StartedAt
		Duration() Duration
	}

	breaksStarted struct {
		startedAt StartedAt
		duration  Duration
	}
)

func NewBreaksStarted(startedAt StartedAt, duration Duration) BreaksStarted {
	return &breaksStarted{
		startedAt,
		duration,
	}
}

func (e *breaksStarted) Type() event.EventType {
	return event.EventTypeBreaksStarted
}

func (e *breaksStarted) StartedAt() StartedAt {
	return e.startedAt
}

func (e *breaksStarted) Duration() Duration {
	return e.duration
}

func (e *breaksStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string `json:"type"`
		StartedAt string `json:"started_at"`
		Duration  string `json:"duration"`
	}{
		e.Type().String(),
		time.Unix(0, e.startedAt.Value()).String(),
		time.Duration(e.duration.Value()).String(),
	})
}

func (e *breaksStarted) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
