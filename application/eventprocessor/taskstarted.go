package eventprocessor

import (
	"github.com/hrfmmr/lyco/domain/entry"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type TaskStartedEventProcessor struct {
	entryRepository entry.Repository
}

func NewTaskStartedEventProcessor(
	entryRepository entry.Repository,
) *TaskStartedEventProcessor {
	return &TaskStartedEventProcessor{
		entryRepository,
	}
}

func (s *TaskStartedEventProcessor) EventType() event.EventType {
	return event.EventTypeTaskStarted
}

func (s *TaskStartedEventProcessor) HandleEvent(e event.Event) {
	ev, ok := e.(task.TaskStarted)
	if !ok {
		logrus.Errorf("❗[TaskStartedEventProcessor] got unexpected event:%T, expecting: task.TaskStarted", e)
		return
	}
	name, err := entry.NewName(ev.Name().Value())
	if err != nil {
		logrus.Errorf("❗[TaskStartedEventProcessor] err:%v", err)
		return
	}
	startedAt, err := entry.NewStartedAt(ev.StartedAt().Value())
	if err != nil {
		logrus.Errorf("❗[TaskStartedEventProcessor] err:%v", err)
		return
	}
	s.entryRepository.Add(entry.NewEntry(name, startedAt))
}
