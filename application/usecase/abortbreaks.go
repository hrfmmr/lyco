package usecase

import (
	"errors"

	"github.com/hrfmmr/lyco/application/appstate"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/sirupsen/logrus"
)

type AbortBreaksUseCase struct {
	appState       appstate.AppState
	taskRepository task.Repository
	pomodorotimer  timer.Timer
}

func NewAbortBreaksUseCase(
	appState appstate.AppState,
	taskRepository task.Repository,
	pomodorotimer timer.Timer,
) *AbortBreaksUseCase {
	return &AbortBreaksUseCase{
		appState,
		taskRepository,
		pomodorotimer,
	}
}

func (u *AbortBreaksUseCase) Execute(arg interface{}) error {
	b := u.appState.CurrentBreaks()
	if b == nil {
		return errors.New("❗[AbortBreaksUseCase] breaks is nil...")
	}
	if b.IsStopped() {
		return nil
	}
	b.Stop()
	u.appState.SetBreaks(nil)
	u.pomodorotimer.Stop()
	t := u.taskRepository.GetCurrent()
	logrus.Infof("🐛current task:%v", t)
	t = task.NewTask(t.Name(), t.Duration())
	logrus.Infof("🐛new task:%v", t)
	u.taskRepository.Save(t)
	return nil
}
