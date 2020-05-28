package usecase

import (
	"errors"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type SwitchTaskUseCase struct {
	taskRepo task.Repository
}

func NewSwitchTaskUseCase(taskRepo task.Repository) *SwitchTaskUseCase {
	return &SwitchTaskUseCase{
		taskRepo,
	}
}

func (u *SwitchTaskUseCase) Execute(arg interface{}) error {
	name, ok := arg.(task.Name)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Name`")
	}
	logrus.Infof("🐛SwitchTaskUseCase#Execute name:%v", name)
	task, err := task.SwitchTask(u.taskRepo.GetCurrent(), name)
	if err != nil {
		return err
	}
	u.taskRepo.Save(task)
	return nil
}
