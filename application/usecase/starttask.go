package usecase

import (
	"errors"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type StartTaskUseCase struct {
	taskRepo task.Repository
}

func NewStartTaskUseCase(taskRepo task.Repository) *StartTaskUseCase {
	return &StartTaskUseCase{
		taskRepo,
	}
}

func (u *StartTaskUseCase) Execute(arg interface{}) error {
	name, ok := arg.(task.Name)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Name`")
	}
	logrus.Infof("🐛StartTaskUseCase#Execute name:%v", name)
	d, err := task.NewDuration(int64(task.DefaultDuration))
	if err != nil {
		return err
	}
	task := task.NewTask(name, d)
	if err := task.Start(time.Now()); err != nil {
		return err
	}
	u.taskRepo.Save(task)
	return nil
}
