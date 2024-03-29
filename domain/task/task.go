package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
)

const (
	DefaultDuration = 25 * time.Minute
)

//go:generate stringer -type=AvailableAction
type AvailableAction int

const (
	AvailableActionStart AvailableAction = iota
	AvailableActionPause
	AvailableActionResume
	AvailableActionStop
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
		Finish() error
		CanStart() bool
		CanPause() bool
		CanResume() bool
		CanStop() bool
		CanFinish() bool
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

func NewTaskWithValues(
	name Name,
	duration Duration,
	startedAt StartedAt,
	elapsed Elapsed,
	status Status,
) Task {
	return &task{name, duration, startedAt, elapsed, status}
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
		return errors.New(fmt.Sprintf("❗cannot start task:%v", t))
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
		return errors.New(fmt.Sprintf("❗cannot pause task:%v", t))
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
		t.status,
	))
	return nil
}

func (t *task) Resume(at time.Time) error {
	if !t.CanResume() {
		return errors.New(fmt.Sprintf("❗cannot resume task:%v", t))
	}
	if err := t.status.Update(NewStatus(TaskStatusRunning)); err != nil {
		return err
	}
	startedAt, err := NewStartedAt(at.UnixNano())
	if err != nil {
		return err
	}
	t.startedAt = startedAt
	event.DefaultPublisher.Publish(NewTaskResumed(
		t.name,
		t.startedAt,
		t.duration,
		t.elapsed,
		t.status,
	))
	return nil
}

func (t *task) Stop() error {
	if !t.CanStop() {
		return errors.New(fmt.Sprintf("❗cannot stop task:%v", t))
	}
	if err := t.status.Update(NewStatus(TaskStatusStopped)); err != nil {
		return err
	}
	now := time.Now().UnixNano()
	elapsed, err := NewElapsed(t.elapsed.Value() + now - t.startedAt.Value())
	if err != nil {
		return err
	}
	t.elapsed = elapsed
	event.DefaultPublisher.Publish(NewTaskStopped(
		t.name,
		t.startedAt,
		t.duration,
		t.elapsed,
		t.status,
	))
	return nil
}

func (t *task) Finish() error {
	if !t.CanFinish() {
		return errors.New(fmt.Sprintf("❗cannot finish task:%v", t))
	}
	if err := t.status.Update(NewStatus(TaskStatusFinished)); err != nil {
		return err
	}
	now := time.Now().UnixNano()
	elapsed, err := NewElapsed(t.elapsed.Value() + now - t.startedAt.Value())
	if err != nil {
		return err
	}
	t.elapsed = elapsed
	event.DefaultPublisher.Publish(NewTaskFinished(
		t.name,
		t.startedAt,
		t.duration,
		t.elapsed,
		t.status,
	))
	return nil
}

func (t *task) AvailableActions() []AvailableAction {
	switch t.Status().Value() {
	case TaskStatusNone, TaskStatusFinished, TaskStatusStopped:
		return []AvailableAction{
			AvailableActionStart,
			AvailableActionSwitch,
		}
	case TaskStatusRunning:
		return []AvailableAction{
			AvailableActionPause,
			AvailableActionStop,
			AvailableActionSwitch,
		}
	case TaskStatusPaused:
		return []AvailableAction{
			AvailableActionResume,
			AvailableActionStop,
			AvailableActionSwitch,
		}
	default:
		return []AvailableAction{
			AvailableActionSwitch,
		}
	}
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

func (t *task) CanStop() bool {
	return t.hasAvailableAction(AvailableActionStop)
}

func (t *task) CanFinish() bool {
	now := time.Now().UnixNano()
	elapsed := t.elapsed.Value() + now - t.startedAt.Value()
	return elapsed >= t.duration.Value()
}

func (t *task) MarshalJSON() ([]byte, error) {
	var startedAt *string
	if t.startedAt != nil {
		startedAtTime := time.Unix(0, t.startedAt.Value()).String()
		startedAt = &startedAtTime
	}
	return json.Marshal(struct {
		Name             string  `json:"name"`
		Duration         string  `json:"duration"`
		StartedAt        *string `json:"started_at"`
		Elapsed          string  `json:"elapsed"`
		Status           string  `json:"status"`
		AvailableActions string  `json:"available_actions"`
	}{
		t.name.Value(),
		time.Duration(t.duration.Value()).String(),
		startedAt,
		time.Duration(t.elapsed.Value()).String(),
		string(t.status.Value()),
		fmt.Sprintf("%v", t.AvailableActions()),
	})
}

func (t *task) String() string {
	b, _ := json.Marshal(t)
	return string(b)
}

func (t *task) hasAvailableAction(action AvailableAction) bool {
	for _, v := range t.AvailableActions() {
		if v == action {
			return true
		}
	}
	return false
}
