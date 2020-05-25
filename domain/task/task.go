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
		Elapsed() time.Duration
		Status() Status
		// behaviors
		Start(at time.Time)
		Pause()
		Resume(at time.Time)
		Stop()
	}

	task struct {
		name      Name
		duration  time.Duration
		startedAt StartedAt
		elapsed   time.Duration
		status    Status
	}
)

func NewTask(name Name, d time.Duration) Task {
	return &task{
		name:     name,
		duration: d,
		status:   NewStatus(TaskStatusNone),
	}
}

func (t *task) Name() Name {
	return t.name
}

func (t *task) Duration() time.Duration {
	return t.duration
}

func (t *task) Elapsed() time.Duration {
	return t.elapsed
}

func (t *task) StartedAt() StartedAt {
	return t.startedAt
}

func (t *task) Status() Status {
	return t.status
}

func (t *task) Start(at time.Time) {
	t.status.Update(NewStatus(TaskStatusRunning))
	t.startedAt = NewStartedAt(at.UnixNano())
	t.elapsed = 0
}

func (t *task) Pause() {
	t.status.Update(NewStatus(TaskStatusPaused))
	now := time.Now().UnixNano()
	t.elapsed += time.Duration(now - t.startedAt.Value())
}

func (t *task) Resume(at time.Time) {
	t.status.Update(NewStatus(TaskStatusRunning))
	t.startedAt = NewStartedAt(at.UnixNano())
}

func (t *task) Stop() {
	t.status.Update(NewStatus(TaskStatusAborted))
	now := time.Now().UnixNano()
	t.elapsed += time.Duration(now - t.startedAt.Value())
}
