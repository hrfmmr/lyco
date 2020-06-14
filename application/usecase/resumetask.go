package usecase

import (
	"time"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type ResumeTaskUseCase struct {
	pomodorotimer  timer.Timer
	taskRepository task.Repository
}

func NewResumeTaskUseCase(pomodorotimer timer.Timer, taskRepository task.Repository) *ResumeTaskUseCase {
	return &ResumeTaskUseCase{
		pomodorotimer,
		taskRepository,
	}
}

func (u *ResumeTaskUseCase) Execute(arg interface{}) error {
	t := u.taskRepository.GetCurrent()
	if err := t.Resume(time.Now()); err != nil {
		return err
	}
	t.Resume(time.Now())
	u.taskRepository.Save(t)
	d, err := timer.NewDuration(t.Duration().Value() - t.Elapsed().Value())
	if err != nil {
		return err
	}
	u.pomodorotimer.Start(timer.TimerModeTask, d)
	return nil
}
