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
		Status() Status
	}

	taskPaused struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
		status    Status
	}
)

func NewTaskPaused(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed, status Status) TaskPaused {
	return &taskPaused{
		name,
		startedAt,
		duration,
		elapsed,
		status,
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

func (e *taskPaused) Status() Status {
	return e.status
}

func (e *taskPaused) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type      string `json:"type"`
		Name      string `json:"name"`
		StartedAt string `json:"started_at"`
		Duration  string `json:"duration"`
		Elapsed   string `json:"elapsed"`
		Status    string `json:"status"`
	}{
		e.Type().String(),
		e.name.Value(),
		time.Unix(0, e.startedAt.Value()).String(),
		time.Duration(e.duration.Value()).String(),
		time.Duration(e.elapsed.Value()).String(),
		string(e.status.Value()),
	})
}

func (e *taskPaused) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
