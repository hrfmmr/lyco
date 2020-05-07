package dto

import (
	"time"

	"github.com/hrfmmr/lyco/domain/task"
)

type (
	TaskDTO interface {
		Name() string
		Duration() time.Duration
		StartedAt() int64
	}

	taskDTO struct {
		name      string
		duration  time.Duration
		startedAt int64
	}
)

func NewTaskDTO() TaskDTO {
	return &taskDTO{}
}

func (t *taskDTO) Name() string {
	return t.name
}

func (t *taskDTO) Duration() time.Duration {
	return t.duration
}

func (t *taskDTO) StartedAt() int64 {
	return t.startedAt
}

func ConvertTaskToDTO(t task.Task) TaskDTO {
	return &taskDTO{
		t.Name().Value(),
		t.Duration(),
		t.StartedAt().Value(),
	}
}
