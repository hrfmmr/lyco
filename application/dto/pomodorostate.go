package dto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

//go:generate stringer -type=PomodoroMode

type PomodoroMode int

const (
	PomodoroModeTask PomodoroMode = iota
	PomodoroModeBreaks
)

//go:generate stringer -type=AvailableTaskAction
type AvailableTaskAction int

const (
	AvailableTaskActionStart AvailableTaskAction = iota
	AvailableTaskActionPause
	AvailableTaskActionResume
	AvailableTaskActionStop
	AvailableTaskActionSwitch
	AvailableTaskActionAbortBreaks
)

type (
	PomodoroState interface {
		Mode() PomodoroMode
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

func NewInitialPomodoroState() PomodoroState {
	return &initialstate{}
}

func (s *initialstate) Mode() PomodoroMode {
	return PomodoroModeTask
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

func NewPomodoroStateWithTask(t task.Task) PomodoroState {
	return &tstate{t}
}

func (s *tstate) Mode() PomodoroMode {
	return PomodoroModeTask
}

func (s *tstate) TaskName() string {
	return s.t.Name().Value()
}

func (s *tstate) RemainsTimerText() string {
	return remainsTimerText(s.remainsDuration())
}

func (s *tstate) remainsDuration() int64 {
	duration, elapsed := s.t.Duration().Value(), s.t.Elapsed().Value()
	switch s.t.Status().Value() {
	case task.TaskStatusPaused:
		return duration - elapsed
	case task.TaskStatusNone, task.TaskStatusStopped:
		return duration
	default:
		if s.t.StartedAt() == nil {
			logrus.Errorf("❗[remainsDuration] startedAt is nil for task:%v", s.t)
		}
		startedAt := s.t.StartedAt().Value()
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
		case task.AvailableActionResume:
			actions = append(actions, AvailableTaskActionResume)
		case task.AvailableActionStop:
			actions = append(actions, AvailableTaskActionStop)
		case task.AvailableActionSwitch:
			actions = append(actions, AvailableTaskActionSwitch)
		}
	}
	return actions
}

func (s *tstate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TaskName         string `json:"task_name"`
		RemainsTimerText string `json:"remains_timer_text"`
		AvailableActions string `json:"available_actions"`
	}{
		s.TaskName(),
		s.RemainsTimerText(),
		fmt.Sprintf("%v", s.AvailableActions()),
	})
}

func (s *tstate) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

//===============================================================

func NewPomodoroStateWithBreaks(b breaks.Breaks) PomodoroState {
	return &bstate{b}
}

func (s *bstate) Mode() PomodoroMode {
	return PomodoroModeBreaks
}

func (s *bstate) TaskName() string {
	return "☕"
}

func (s *bstate) RemainsTimerText() string {
	return remainsTimerText(s.remainsDuration())
}

func (s *bstate) remainsDuration() int64 {
	if s.b.StartedAt() == nil {
		logrus.Errorf("❗Breaks seem not to have been started...")
		return s.b.Duration().Value()
	}
	to := s.b.StartedAt().Value() + s.b.Duration().Value()
	now := time.Now().UnixNano()
	return to - now
}

func (s *bstate) AvailableActions() []AvailableTaskAction {
	return []AvailableTaskAction{
		AvailableTaskActionAbortBreaks,
	}
}

func (s *bstate) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TaskName         string `json:"task_name"`
		RemainsTimerText string `json:"remains_timer_text"`
		AvailableActions string `json:"available_actions"`
	}{
		s.TaskName(),
		s.RemainsTimerText(),
		fmt.Sprintf("%v", s.AvailableActions()),
	})
}

func (s *bstate) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

//===============================================================

func remainsTimerText(remainsDurationNano int64) string {
	rsec := remainsDurationNano / 1e9
	return fmt.Sprintf("%02d:%02d", int(rsec/60)%60, rsec%60)
}
