package task

import (
	"time"
)

const (
	DefaultDuration = 25 * time.Minute
)

type (
	Task interface {
		// props
		Name() Name
		Duration() time.Duration
		StartedAt() StartedAt
		Status() Status
		// behaviors
		Start(at time.Time)
		Stop()
		Pause()
		Resume()
	}

	task struct {
		name      Name
		duration  time.Duration
		startedAt StartedAt
		status    Status
	}
)

func NewTask(name Name, d time.Duration) Task {
	return &task{
		name:     name,
		duration: d,
		status:   TaskStatusNone,
	}
}

func (t *task) Name() Name {
	return t.name
}

func (t *task) Duration() time.Duration {
	return t.duration
}

func (t *task) StartedAt() StartedAt {
	return t.startedAt
}

func (t *task) Status() Status {
	return t.status
}

func (t *task) Start(at time.Time) {
	t.startedAt = NewStartedAt(at.UnixNano())
	t.status = TaskStatusRunning
}
func (t *task) Stop()   {}
func (t *task) Pause()  {}
func (t *task) Resume() {}
