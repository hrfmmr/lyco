package task

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

type (
	TaskResumed interface {
		event.Event
		Name() Name
		StartedAt() StartedAt
		Duration() Duration
		Elapsed() Elapsed
	}

	taskResumed struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
	}
)

func NewTaskResumed(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed) TaskResumed {
	return &taskResumed{
		name,
		startedAt,
		duration,
		elapsed,
	}
}

func (e *taskResumed) Type() event.EventType {
	return event.EventTypeTaskResumed
}

func (e *taskResumed) Name() Name {
	return e.name
}
func (e *taskResumed) StartedAt() StartedAt {
	return e.startedAt
}
func (e *taskResumed) Duration() Duration {
	return e.duration
}
func (e *taskResumed) Elapsed() Elapsed {
	return e.elapsed
}

func (e *taskResumed) MarshalJSON() ([]byte, error) {
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

func (e *taskResumed) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
