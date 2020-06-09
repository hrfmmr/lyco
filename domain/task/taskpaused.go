package task

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

type (
	TaskPaused interface {
		event.Event
		Name() Name
		StartedAt() StartedAt
		Duration() Duration
		Elapsed() Elapsed
	}

	taskPaused struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
	}
)

func NewTaskPaused(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed) TaskPaused {
	return &taskPaused{
		name,
		startedAt,
		duration,
		elapsed,
	}
}

func (e *taskPaused) Type() event.EventType {
	return event.EventTypeTaskPaused
}

func (e *taskPaused) Name() Name {
	return e.name
}
func (e *taskPaused) StartedAt() StartedAt {
	return e.startedAt
}
func (e *taskPaused) Duration() Duration {
	return e.duration
}
func (e *taskPaused) Elapsed() Elapsed {
	return e.elapsed
}

func (e *taskPaused) MarshalJSON() ([]byte, error) {
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

func (e *taskPaused) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
