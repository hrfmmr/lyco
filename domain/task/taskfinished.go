package task

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

type (
	TaskFinished interface {
		event.Event
		Name() Name
		StartedAt() StartedAt
		Duration() Duration
		Elapsed() Elapsed
		Status() Status
	}

	taskFinished struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
		status    Status
	}
)

func NewTaskFinished(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed, status Status) TaskFinished {
	return &taskFinished{
		name,
		startedAt,
		duration,
		elapsed,
		status,
	}
}

func (e *taskFinished) Type() event.EventType {
	return event.EventTypeTaskFinished
}

func (e *taskFinished) Name() Name {
	return e.name
}

func (e *taskFinished) StartedAt() StartedAt {
	return e.startedAt
}

func (e *taskFinished) Duration() Duration {
	return e.duration
}

func (e *taskFinished) Elapsed() Elapsed {
	return e.elapsed
}

func (e *taskFinished) Status() Status {
	return e.status
}

func (e *taskFinished) MarshalJSON() ([]byte, error) {
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

func (e *taskFinished) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
