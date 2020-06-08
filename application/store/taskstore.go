package store

import (
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type TaskStore interface {
	Store
	GetState() dto.TaskState
	SetState(state dto.TaskState)
}

type taskStore struct {
	onChangeCh     chan Store
	taskRepository task.Repository
	state          dto.TaskState
}

func NewTaskStore(taskRepository task.Repository) TaskStore {
	return &taskStore{
		make(chan Store, 1),
		taskRepository,
		dto.NewInitialTaskState(),
	}
}

func (s *taskStore) RecvPayload(p usecase.Payload) {
	logrus.Infof("🐛taskStore#RecvPayload p:%v", p)
	if t := s.taskRepository.GetCurrent(); t != nil {
		newstate := dto.NewTaskStateWithTask(t)
		s.state = newstate
	}
}
func (s *taskStore) OnChange() <-chan Store {
	return s.onChangeCh
}

func (s *taskStore) GetState() dto.TaskState {
	return s.state
}

func (s *taskStore) SetState(state dto.TaskState) {
	s.state = state
	s.onChangeCh <- s
}
