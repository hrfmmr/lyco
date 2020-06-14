package eventprocessor

import (
	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/sirupsen/logrus"
)

type TimerTickedEventProcessor struct {
	appContext     application.AppContext
	taskRepository task.Repository
}

func NewTimerTickedEventProcessor(
	appContext application.AppContext,
	taskRepository task.Repository,
) *TimerTickedEventProcessor {
	return &TimerTickedEventProcessor{
		appContext,
		taskRepository,
	}
}

func (s *TimerTickedEventProcessor) EventType() event.EventType {
	return event.EventTypeTimerTicked
}

func (s *TimerTickedEventProcessor) HandleEvent(e event.Event) {
	ev, ok := e.(timer.TimerTicked)
	if !ok {
		logrus.Errorf("❗[TimerStateUpdater] got unexpected event:%T, expecting: task.TimerTicked", e)
		return
	}
	switch ev.Mode() {
	case timer.TimerModeTask:
		t := s.taskRepository.GetCurrent()
		newstate := dto.NewPomodoroStateWithTask(t)
		s.appContext.StoreGroup().TaskStore().SetState(newstate)
	case timer.TimerModeBreaks:
		b := s.appContext.AppState().CurrentBreaks()
		if b == nil {
			logrus.Errorf("❗[TimerStateUpdater] breaks is nil...")
			return
		}
		newstate := dto.NewPomodoroStateWithBreaks(b)
		s.appContext.StoreGroup().TaskStore().SetState(newstate)
	}
}
