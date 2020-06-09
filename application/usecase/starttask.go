package usecase

import (
	"errors"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/sirupsen/logrus"
)

type StartTaskUseCase struct {
	pomodorotimer  timer.Timer
	taskRepository task.Repository
}

func NewStartTaskUseCase(pomodorotimer timer.Timer, taskRepository task.Repository) *StartTaskUseCase {
	return &StartTaskUseCase{
		pomodorotimer,
		taskRepository,
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
	u.taskRepository.Save(task)
	u.pomodorotimer.Start(task.Duration(), task.Elapsed())
	return nil
}
