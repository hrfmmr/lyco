package usecase

import (
	"errors"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type StopTaskUseCase struct {
	taskRepo task.Repository
}

func NewStopTaskUseCase(taskRepo task.Repository) *StopTaskUseCase {
	return &StopTaskUseCase{
		taskRepo,
	}
}

func (u *StopTaskUseCase) Execute(arg interface{}) error {
	task, ok := arg.(task.Task)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Task`")
	}
	logrus.Infof("🐛StopTaskUseCase#Execute task:%v", task)
	task.Stop()
	u.taskRepo.Save(task)
	return nil
}
