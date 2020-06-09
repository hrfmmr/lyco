package usecase

import (
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type PauseTaskUseCase struct {
	pomodorotimer  timer.Timer
	taskRepository task.Repository
}

func NewPauseTaskUseCase(pomodorotimer timer.Timer, taskRepository task.Repository) *PauseTaskUseCase {
	return &PauseTaskUseCase{
		pomodorotimer,
		taskRepository,
	}
}

func (u *PauseTaskUseCase) Execute(arg interface{}) error {
	t := u.taskRepository.GetCurrent()
	if !t.CanPause() {
		// ignore
		return nil
	}
	t.Pause()
	u.taskRepository.Save(t)
	u.pomodorotimer.Stop()
	return nil
}
