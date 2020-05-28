package usecase

import (
	"errors"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type AbortBreaksUseCase struct {
	taskRepo task.Repository
}

func NewAbortBreaksUseCase(taskRepo task.Repository) *AbortBreaksUseCase {
	return &AbortBreaksUseCase{
		taskRepo,
	}
}

func (u *AbortBreaksUseCase) Execute(arg interface{}) error {
	name, ok := arg.(task.Name)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Name`")
	}
	logrus.Infof("🐛AbortBreaksUseCase#Execute name:%v", name)
	task := task.NewTask(name, task.DefaultDuration)
	u.taskRepo.Save(task)
	return nil
}
