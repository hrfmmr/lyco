package eventprocessor

import (
	"fmt"

	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type TimerStateUpdater struct {
	taskStore      store.TaskStore
	taskRepository task.Repository
}

func NewTimerStateUpdater(taskStore store.TaskStore, taskRepository task.Repository) *TimerStateUpdater {
	return &TimerStateUpdater{
		taskStore,
		taskRepository,
	}
}

func (s *TimerStateUpdater) EventType() event.EventType {
	return event.EventTypeTimerTicked
}

func (s *TimerStateUpdater) HandleEvent(e event.Event) {
	_, ok := e.(timer.TimerTicked)
	if !ok {
		panic(fmt.Sprintf("😕 got unexpected event:%v, expecting: task.TimerTicked", e))
	}
	t := s.taskRepository.GetCurrent()
	newstate := dto.NewTaskStateWithTask(t)
	s.taskStore.SetState(newstate)
}
