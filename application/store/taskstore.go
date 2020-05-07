package store

import (
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type TaskStore interface {
	Store
	GetState() TaskState
}

type taskStore struct {
	taskRepo task.Repository
	state    TaskState
}

func NewTaskStore(taskRepo task.Repository) TaskStore {
	return &taskStore{
		taskRepo: taskRepo,
		state:    NewTaskState(),
	}
}

func (s *taskStore) RecvPayload(p usecase.Payload) {
	logrus.Infof("🐛taskStore#RecvPayload p:%v", p)
	if task := s.taskRepo.GetCurrent(); task != nil {
		s.state.Update(task)
	}
}

func (s *taskStore) GetState() TaskState {
	return s.state
}
