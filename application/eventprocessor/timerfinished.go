package eventprocessor

import (
	"time"

	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/hrfmmr/lyco/utils/notifier"
	log "github.com/sirupsen/logrus"
)

//TODO: make configurable
const longBreaksPerPoms = 4

type TimerFinishedEventProcessor struct {
	appContext     application.AppContext
	taskRepository task.Repository
	pomodorotimer  timer.Timer
	taskFinCount   uint
}

func NewTimerFinishedEventProcessor(
	appContext application.AppContext,
	taskRepository task.Repository,
	pomodorotimer timer.Timer,
) *TimerFinishedEventProcessor {
	return &TimerFinishedEventProcessor{
		appContext,
		taskRepository,
		pomodorotimer,
		0,
	}
}

func (p *TimerFinishedEventProcessor) EventType() event.EventType {
	return event.EventTypeTimerFinished
}

func (p *TimerFinishedEventProcessor) HandleEvent(e event.Event) {
	ev, ok := e.(timer.TimerFinished)
	if !ok {
		log.Errorf("❗[TimerFinishedEventProcessor] got unexpected event:%T, expecting: task.TimerFinished", e)
		return
	}
	switch ev.Mode() {
	case timer.TimerModeTask:
		//TODO: responsibility segregation
		t := p.taskRepository.GetCurrent()
		if t.CanFinish() {
			t.Finish()
		} else {
			log.Warnf("⚠ task:%v should be able to finish", t)
		}
		p.taskRepository.Save(t)

		p.taskFinCount++
		var b breaks.Breaks
		if p.taskFinCount%longBreaksPerPoms == 0 {
			b = breaks.LongDefault()
		} else {
			b = breaks.ShortDefault()
		}
		if err := b.Start(); err != nil {
			log.Error(err)
			return
		}
		d, err := timer.NewDuration(b.Duration().Value())
		if err != nil {
			log.Error(err)
			return
		}
		p.pomodorotimer.Start(timer.TimerModeBreaks, d)
		p.appContext.AppState().SetBreaks(b)
		newstate := dto.NewPomodoroStateWithBreaks(b)
		p.appContext.StoreGroup().TaskStore().SetState(newstate)
		notifier.NotifyForBreaksStart(notifier.New(), b)
	case timer.TimerModeBreaks:
		//TODO: responsibility segregation
		t := p.taskRepository.GetCurrent()
		t = task.NewTask(t.Name(), t.Duration())
		if err := t.Start(time.Now()); err != nil {
			log.Error(err)
			return
		}
		p.taskRepository.Save(t)
		d, err := timer.NewDuration(t.Duration().Value())
		if err != nil {
			log.Error(err)
			return
		}
		p.pomodorotimer.Start(timer.TimerModeTask, d)
		p.appContext.AppState().SetBreaks(nil)
		newstate := dto.NewPomodoroStateWithTask(t)
		p.appContext.StoreGroup().TaskStore().SetState(newstate)
		notifier.NotifyForBreaksEnd(notifier.New(), t)
	}
}
