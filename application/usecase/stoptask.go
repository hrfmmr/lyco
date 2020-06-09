package usecase

import (
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type StopTaskUseCase struct {
	pomodorotimer  timer.Timer
	taskRepository task.Repository
}

func NewStopTaskUseCase(pomodorotimer timer.Timer, taskRepository task.Repository) *StopTaskUseCase {
	return &StopTaskUseCase{
		pomodorotimer,
		taskRepository,
	}
}

func (u *StopTaskUseCase) Execute(arg interface{}) error {
	t := u.taskRepository.GetCurrent()
	if !t.CanStop() {
		// ignore
		return nil
	}
	t.Stop()
	u.taskRepository.Save(t)
	u.pomodorotimer.Stop()
	return nil
}
