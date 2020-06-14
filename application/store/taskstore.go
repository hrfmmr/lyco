package store

import (
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type TaskStore interface {
	Store
	GetState() dto.PomodoroState
	SetState(state dto.PomodoroState)
}

type taskStore struct {
	onChangeCh     chan Store
	taskRepository task.Repository
	state          dto.PomodoroState
}

func NewTaskStore(taskRepository task.Repository) TaskStore {
	return &taskStore{
		make(chan Store, 1),
		taskRepository,
		dto.NewInitialPomodoroState(),
	}
}

func (s *taskStore) RecvPayload(p usecase.Payload) {
	logrus.Infof("🐛taskStore#RecvPayload p:%v", p)
	if t := s.taskRepository.GetCurrent(); t != nil {
		newstate := dto.NewPomodoroStateWithTask(t)
		s.state = newstate
	}
}
func (s *taskStore) OnChange() <-chan Store {
	return s.onChangeCh
}

func (s *taskStore) GetState() dto.PomodoroState {
	return s.state
}

func (s *taskStore) SetState(state dto.PomodoroState) {
	s.state = state
	s.onChangeCh <- s
}
