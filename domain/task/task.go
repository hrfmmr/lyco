package task

import (
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
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
		Duration() Duration
		StartedAt() StartedAt
		Elapsed() Elapsed
		Status() Status
		AvailableActions() []AvailableAction
		// behaviors
		Start(at time.Time) error
		Pause() error
		Resume(at time.Time) error
		Stop() error
		CanStart() bool
		CanPause() bool
		CanResume() bool
		CanAbort() bool
	}

	task struct {
		name      Name
		duration  Duration
		startedAt StartedAt
		elapsed   Elapsed
		status    Status
	}
)

func NewTask(name Name, d Duration) Task {
	elapsed, _ := NewElapsed(0)
	return &task{
		name:     name,
		duration: d,
		elapsed:  elapsed,
		status:   NewStatus(TaskStatusNone),
	}
}

func NewTaskWithElapsed(name Name, d Duration, elapsed Elapsed) Task {
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

func (t *task) Duration() Duration {
	return t.duration
}

func (t *task) Elapsed() Elapsed {
	return t.elapsed
}

func (t *task) StartedAt() StartedAt {
	return t.startedAt
}

func (t *task) Status() Status {
	return t.status
}

func (t *task) Start(at time.Time) error {
	if !t.CanStart() {
		return NewInvalidStatusTransition(fmt.Sprintf("❗Can't Start from task status:%v", t.status.Value()))
	}
	if err := t.status.Update(NewStatus(TaskStatusRunning)); err != nil {
		return err
	}
	startedAt, err := NewStartedAt(at.UnixNano())
	if err != nil {
		return err
	}
	t.startedAt = startedAt
	event.DefaultPublisher.Publish(NewTaskStarted(
		t.name,
		t.startedAt,
		t.duration,
		t.elapsed,
	))
	return nil
}

func (t *task) Pause() error {
	if !t.CanPause() {
		return NewInvalidStatusTransition(fmt.Sprintf("❗Can't Pause from task status:%v", t.status.Value()))
	}
	if err := t.status.Update(NewStatus(TaskStatusPaused)); err != nil {
		return err
	}
	now := time.Now().UnixNano()
	elapsed, err := NewElapsed(t.elapsed.Value() + now - t.startedAt.Value())
	if err != nil {
		return err
	}
	t.elapsed = elapsed
	event.DefaultPublisher.Publish(NewTaskPaused(
		t.name,
		t.startedAt,
		t.duration,
		t.elapsed,
	))
	return nil
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

func (t *task) Stop() error {
	t.status.Update(NewStatus(TaskStatusAborted))
	now := time.Now().UnixNano()
	elapsed, err := NewElapsed(t.elapsed.Value() + now - t.startedAt.Value())
	if err != nil {
		return err
	}
	t.elapsed = elapsed
	return nil
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
	d, err := NewDuration(int64(DefaultDuration))
	if err != nil {
		return nil, err
	}
	t := NewTaskWithElapsed(to, d, current.Elapsed())
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
