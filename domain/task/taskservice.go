package task

import (
	"time"
)

type (
	TaskService interface {
		SwitchTask(name Name, duration Duration) (Task, error)
	}

	taskservice struct {
		repository Repository
	}
)

func NewTaskService(repository Repository) TaskService {
	return &taskservice{
		repository,
	}
}

func (s *taskservice) SwitchTask(name Name, duration Duration) (Task, error) {
	current := s.repository.GetCurrent()
	switch current.Status().Value() {
	case TaskStatusRunning, TaskStatusPaused:
		current.Stop()
	}
	t := NewTask(name, duration)
	if err := t.Start(time.Now()); err != nil {
		return nil, err
	}
	return t, nil
}
