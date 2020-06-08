package dto

import (
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type AvailableTaskAction int

const (
	AvailableTaskActionStart = iota
	AvailableTaskActionPause
	AvailableTaskActionResume
	AvailableTaskActionAbort
	AvailableTaskActionSwitch
	AvailableTaskActionAbortBreaks
)

type (
	TaskState interface {
		TaskName() string
		RemainsTimerText() string
		AvailableActions() []AvailableTaskAction
	}

	initialstate struct{}

	tstate struct {
		t task.Task
	}

	bstate struct {
		b breaks.Breaks
	}
)

func NewInitialTaskState() TaskState {
	return &initialstate{}
}

func (s *initialstate) TaskName() string {
	return ""
}

func (s *initialstate) RemainsTimerText() string {
	return ""
}

func (s *initialstate) AvailableActions() []AvailableTaskAction {
	return []AvailableTaskAction{}
}

//===============================================================

func NewTaskStateWithTask(t task.Task) TaskState {
	return &tstate{t}
}

func (s *tstate) TaskName() string {
	return s.t.Name().Value()
}

func (s *tstate) RemainsTimerText() string {
	return remainsTimerText(s.remainsDuration())
}

func (s *tstate) remainsDuration() int64 {
	duration, elapsed, startedAt := s.t.Duration().Value(), s.t.Elapsed().Value(), s.t.StartedAt().Value()
	switch s.t.Status().Value() {
	case task.TaskStatusPaused:
		return duration - elapsed
	case task.TaskStatusNone, task.TaskStatusAborted:
		return duration
	default:
		to := startedAt + (duration - elapsed)
		now := time.Now().UnixNano()
		return to - now
	}
}

func (s *tstate) AvailableActions() []AvailableTaskAction {
	actions := []AvailableTaskAction{}
	for _, v := range s.t.AvailableActions() {
		switch v {
		case task.AvailableActionStart:
			actions = append(actions, AvailableTaskActionStart)
		case task.AvailableActionPause:
			actions = append(actions, AvailableTaskActionPause)
		case task.AvailableActionAbort:
			actions = append(actions, AvailableTaskActionAbort)
		case task.AvailableActionSwitch:
			actions = append(actions, AvailableTaskActionSwitch)
		}
	}
	return actions
}

//===============================================================

func NewTaskStateWithBreaks(b breaks.Breaks) TaskState {
	return &bstate{b}
}

func (s *bstate) TaskName() string {
	return ""
}

func (s *bstate) RemainsTimerText() string {
	return remainsTimerText(s.remainsDuration())
}

func (s *bstate) remainsDuration() int64 {
	if s.b.StartedAt() == nil {
		logrus.Errorf("❗Breaks seem not to have been started...")
		return int64(s.b.Duration())
	}
	to := s.b.StartedAt().Value() + int64(s.b.Duration())
	now := time.Now().UnixNano()
	return to - now
}

func (s *bstate) AvailableActions() []AvailableTaskAction {
	return []AvailableTaskAction{
		AvailableTaskActionAbortBreaks,
	}
}

//===============================================================

func remainsTimerText(remainsDurationNano int64) string {
	rsec := remainsDurationNano / 1e9
	return fmt.Sprintf("%02d:%02d", int(rsec/60)%60, rsec%60)
}