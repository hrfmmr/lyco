package dto

import (
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
)

const (
	TaskStatusNone     = "none"
	TaskStatusRunning  = "running"
	TaskStatusPaused   = "paused"
	TaskStatusAborted  = "aborted"
	TaskStatusFinished = "finished"
)

type (
	TaskDTO interface {
		Name() string
		Duration() int64
		StartedAt() int64
		Elapsed() int64
		Status() string
		RemainsDuration() int64
		RemainsTimerText() string
	}

	taskDTO struct {
		name      string
		duration  int64
		startedAt int64
		elapsed   int64
		status    string
	}
)

func NewTaskDTO() TaskDTO {
	return &taskDTO{
		status: TaskStatusNone,
	}
}

func (t *taskDTO) Name() string {
	return t.name
}

func (t *taskDTO) Duration() int64 {
	return t.duration
}

func (t *taskDTO) StartedAt() int64 {
	return t.startedAt
}

func (t *taskDTO) Elapsed() int64 {
	return t.elapsed
}

func (t *taskDTO) Status() string {
	return t.status
}

func (t *taskDTO) RemainsDuration() int64 {
	duration, elapsed, startedAt := t.Duration(), t.Elapsed(), t.StartedAt()
	switch t.Status() {
	case TaskStatusPaused:
		return duration - elapsed
	default:
		to := startedAt + (duration - elapsed)
		now := time.Now().UnixNano()
		return to - now
	}
}

func (t *taskDTO) RemainsTimerText() string {
	rsec := t.RemainsDuration() / 1e9
	return fmt.Sprintf("%02d:%02d", int(rsec/60)%60, rsec%60)
}

func ConvertTaskToDTO(t task.Task) TaskDTO {
	return &taskDTO{
		t.Name().Value(),
		int64(t.Duration()),
		t.StartedAt().Value(),
		int64(t.Elapsed()),
		string(t.Status()),
	}
}
