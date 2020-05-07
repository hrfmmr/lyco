package store

import (
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type TaskState interface {
	State() dto.TaskDTO
	Update(t task.Task)
}

type taskState struct {
	task dto.TaskDTO
}

func NewTaskState() TaskState {
	return &taskState{
		task: dto.NewTaskDTO(),
	}
}

func (s *taskState) State() dto.TaskDTO {
	return s.task
}

func (s *taskState) Update(t task.Task) {
	logrus.Infof("🐛taskState#Update t:%v", t)
	if t == nil {
		return
	}
	s.task = dto.ConvertTaskToDTO(t)
}
