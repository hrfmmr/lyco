package event

import (
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
)

type Event interface {
	Type() EventType
}

type (
	TaskStarted interface {
		Event
		Name() string
		StartedAt() int64
		Duration() time.Duration
		Elapsed() time.Duration
	}

	taskStarted struct {
		name      string
		startedAt int64
		duration  time.Duration
		elapsed   time.Duration
	}
)

func NewTaskStarted(name string, startedAt int64, duration, elapsed time.Duration) TaskStarted {
	return &taskStarted{
		name,
		startedAt,
		duration,
		elapsed,
	}
}

func (e *taskStarted) Type() EventType {
	return EventTypeTaskStarted
}

func (e *taskStarted) Name() string {
	return e.name
}
func (e *taskStarted) StartedAt() int64 {
	return e.startedAt
}
func (e *taskStarted) Duration() time.Duration {
	return e.duration
}
func (e *taskStarted) Elapsed() time.Duration {
	return e.elapsed
}

func (e *taskStarted) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		StartedAt int64  `json:"started_at"`
		Duration  string `json:"duration"`
		Elapsed   string `json:"elapsed"`
	}{
		e.Type().String(),
		e.name,
		e.startedAt,
		e.duration.String(),
		e.elapsed.String(),
	})
}

func (e *taskStarted) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		logrus.Fatalf("❗Failed json.Marshal err:%v", err)
	}
	return string(b)
}
