package task

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/sirupsen/logrus"
)

type (
	TaskStarted interface {
		event.Event
		Name() Name
		StartedAt() StartedAt
		Duration() Duration
		Elapsed() Elapsed
	}

	taskStarted struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
	}
)

func NewTaskStarted(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed) TaskStarted {
	return &taskStarted{
		name,
		startedAt,
		duration,
		elapsed,
	}
}

func (e *taskStarted) Type() event.EventType {
	return event.EventTypeTaskStarted
}

func (e *taskStarted) Name() Name {
	return e.name
}
func (e *taskStarted) StartedAt() StartedAt {
	return e.startedAt
}
func (e *taskStarted) Duration() Duration {
	return e.duration
}
func (e *taskStarted) Elapsed() Elapsed {
	return e.elapsed
}

func (e *taskStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		StartedAt string `json:"started_at"`
		Duration  string `json:"duration"`
		Elapsed   string `json:"elapsed"`
	}{
		e.Type().String(),
		e.name.Value(),
		time.Unix(0, e.startedAt.Value()).String(),
		time.Duration(e.duration.Value()).String(),
		time.Duration(e.elapsed.Value()).String(),
	})
}

func (e *taskStarted) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
