package usecase

import (
	"errors"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type PauseTaskUseCase struct {
	taskRepo task.Repository
}

func NewPauseTaskUseCase(taskRepo task.Repository) *PauseTaskUseCase {
	return &PauseTaskUseCase{
		taskRepo,
	}
}

func (u *PauseTaskUseCase) Execute(arg interface{}) error {
	task, ok := arg.(task.Task)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Task`")
	}
	logrus.Infof("🐛PauseTaskUseCase#Execute task:%v", task)
	task.Pause()
	u.taskRepo.Save(task)
	return nil
}
