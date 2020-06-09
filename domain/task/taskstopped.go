package task

import (
	"encoding/json"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

type (
	TaskStopped interface {
		event.Event
		Name() Name
		StartedAt() StartedAt
		Duration() Duration
		Elapsed() Elapsed
	}

	taskStopped struct {
		name      Name
		startedAt StartedAt
		duration  Duration
		elapsed   Elapsed
	}
)

func NewTaskStopped(name Name, startedAt StartedAt, duration Duration, elapsed Elapsed) TaskStopped {
	return &taskStopped{
		name,
		startedAt,
		duration,
		elapsed,
	}
}

func (e *taskStopped) Type() event.EventType {
	return event.EventTypeTaskStopped
}

func (e *taskStopped) Name() Name {
	return e.name
}
func (e *taskStopped) StartedAt() StartedAt {
	return e.startedAt
}
func (e *taskStopped) Duration() Duration {
	return e.duration
}
func (e *taskStopped) Elapsed() Elapsed {
	return e.elapsed
}

func (e *taskStopped) MarshalJSON() ([]byte, error) {
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

func (e *taskStopped) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
