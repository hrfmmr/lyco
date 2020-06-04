package usecase

import (
	"errors"
	"time"

	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type StartTaskUseCase struct {
	taskRepo       task.Repository
	eventPublisher event.Publisher
}

func NewStartTaskUseCase(taskRepo task.Repository, eventPublisher event.Publisher) *StartTaskUseCase {
	return &StartTaskUseCase{
		taskRepo,
		eventPublisher,
	}
}

func (u *StartTaskUseCase) Execute(arg interface{}) error {
	name, ok := arg.(task.Name)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Name`")
	}
	logrus.Infof("🐛StartTaskUseCase#Execute name:%v", name)
	task := task.NewTask(name, task.DefaultDuration)
	if err := task.Start(time.Now()); err != nil {
		return err
	}
	u.taskRepo.Save(task)
	u.eventPublisher.Publish(event.NewTaskStarted(
		task.Name().Value(),
		task.StartedAt().Value(),
		task.Duration(),
		task.Elapsed(),
	))
	return nil
}
