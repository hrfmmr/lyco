package task

import (
	"time"
)

const (
	DefaultDuration = 25 * time.Minute
)

type AvailableAction int

const (
	AvailableActionStart = iota
	AvailableActionPause
	AvailableActionResume
	AvailableActionAbort
	AvailableActionSwitch
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
		Start(at time.Time) error
		Pause()
		Resume(at time.Time) error
		Stop()
		// utils
		AvailableActions() []AvailableAction
		CanStart() bool
		CanPause() bool
		CanResume() bool
		CanAbort() bool
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

func NewTaskWithElapsed(name Name, d, elapsed time.Duration) Task {
	return &task{
		name:     name,
		duration: d,
		elapsed:  elapsed,
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

func (t *task) Start(at time.Time) error {
	t.status.Update(NewStatus(TaskStatusRunning))
	startedAt, err := NewStartedAt(at.UnixNano())
	if err != nil {
		return err
	}
	t.startedAt = startedAt
	return nil
}

func (t *task) Pause() {
	t.status.Update(NewStatus(TaskStatusPaused))
	now := time.Now().UnixNano()
	t.elapsed += time.Duration(now - t.startedAt.Value())
}

func (t *task) Resume(at time.Time) error {
	t.status.Update(NewStatus(TaskStatusRunning))
	startedAt, err := NewStartedAt(at.UnixNano())
	if err != nil {
		return err
	}
	t.startedAt = startedAt
	return nil
}

func (t *task) Stop() {
	t.status.Update(NewStatus(TaskStatusAborted))
	now := time.Now().UnixNano()
	t.elapsed += time.Duration(now - t.startedAt.Value())
}

func (t *task) AvailableActions() []AvailableAction {
	switch t.Status().Value() {
	case TaskStatusNone, TaskStatusFinished, TaskStatusAborted:
		return []AvailableAction{
			AvailableActionStart,
			AvailableActionSwitch,
		}
	case TaskStatusRunning:
		return []AvailableAction{
			AvailableActionPause,
			AvailableActionAbort,
			AvailableActionSwitch,
		}
	case TaskStatusPaused:
		return []AvailableAction{
			AvailableActionResume,
			AvailableActionAbort,
			AvailableActionSwitch,
		}
	default:
		return []AvailableAction{
			AvailableActionSwitch,
		}
	}
}

func SwitchTask(current Task, to Name) (Task, error) {
	if current.CanAbort() {
		current.Stop()
	}
	t := NewTaskWithElapsed(to, DefaultDuration, current.Elapsed())
	if err := t.Start(time.Now()); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *task) CanStart() bool {
	return t.hasAvailableAction(AvailableActionStart)
}

func (t *task) CanPause() bool {
	return t.hasAvailableAction(AvailableActionPause)
}

func (t *task) CanResume() bool {
	return t.hasAvailableAction(AvailableActionResume)
}

func (t *task) CanAbort() bool {
	return t.hasAvailableAction(AvailableActionAbort)
}

func (t *task) hasAvailableAction(action AvailableAction) bool {
	for _, v := range t.AvailableActions() {
		if v == action {
			return true
		}
	}
	return false
}
