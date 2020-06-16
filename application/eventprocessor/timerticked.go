package eventprocessor

import (
	"time"

	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/domain/entry"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/sirupsen/logrus"
)

type TimerTickedEventProcessor struct {
	appContext      application.AppContext
	metricsStore    store.MetricsStore
	taskRepository  task.Repository
	entryRepository entry.Repository
}

func NewTimerTickedEventProcessor(
	appContext application.AppContext,
	metricsStore store.MetricsStore,
	taskRepository task.Repository,
	entryRepository entry.Repository,
) *TimerTickedEventProcessor {
	return &TimerTickedEventProcessor{
		appContext,
		metricsStore,
		taskRepository,
		entryRepository,
	}
}

func (p *TimerTickedEventProcessor) EventType() event.EventType {
	return event.EventTypeTimerTicked
}

func (p *TimerTickedEventProcessor) HandleEvent(e event.Event) {
	ev, ok := e.(timer.TimerTicked)
	if !ok {
		logrus.Errorf("❗[TimerTickedEventProcessor] got unexpected event:%T, expecting: task.TimerTicked", e)
		return
	}
	switch ev.Mode() {
	case timer.TimerModeTask:
		t := p.taskRepository.GetCurrent()
		newPomodoroState := dto.NewPomodoroStateWithTask(t)
		p.appContext.StoreGroup().TaskStore().SetState(newPomodoroState)

		entr, err := p.entryRepository.GetLatest()
		if err != nil {
			logrus.Errorf("❗[TimerTickedEventProcessor] err:%v", err)
			return
		}
		elapsed, err := entry.NewElapsed(int64(time.Second))
		if err != nil {
			logrus.Errorf("❗[TimerTickedEventProcessor] err:%v", err)
			return
		}
		entr.IncrementElapsed(elapsed)
		p.entryRepository.Update(entr)
		entries, err := p.entryRepository.GetAll()
		if err != nil {
			logrus.Errorf("❗[TimerTickedEventProcessor] err:%v", err)
			return
		}
		newMetricsState := dto.NewMetricsState(entries, t.Duration().Value())
		p.metricsStore.SetState(newMetricsState)
	case timer.TimerModeBreaks:
		b := p.appContext.AppState().CurrentBreaks()
		if b == nil {
			logrus.Errorf("❗[TimerTickedEventProcessor] breaks is nil...")
			return
		}
		newstate := dto.NewPomodoroStateWithBreaks(b)
		p.appContext.StoreGroup().TaskStore().SetState(newstate)
	}
}
