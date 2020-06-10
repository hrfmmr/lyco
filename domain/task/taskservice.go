package task

import (
	"errors"
	"fmt"
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
	var t Task
	switch current.Status().Value() {
	case TaskStatusNone, TaskStatusFinished, TaskStatusStopped:
		t = NewTask(name, duration)
	case TaskStatusRunning:
		current.Stop()
		t = NewTaskWithElapsed(name, duration, current.Elapsed())
	case TaskStatusPaused:
		t = NewTaskWithElapsed(name, duration, current.Elapsed())
		current.Stop()
	default:
		return nil, errors.New(fmt.Sprintf("❗[SwitchTask]unexpected status:%v", current.Status().Value()))
	}
	if err := t.Start(time.Now()); err != nil {
		return nil, err
	}
	return t, nil
}
