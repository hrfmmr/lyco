package usecase

import (
	"errors"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type ResumeTaskUseCase struct {
	taskRepo task.Repository
}

func NewResumeTaskUseCase(taskRepo task.Repository) *ResumeTaskUseCase {
	return &ResumeTaskUseCase{
		taskRepo,
	}
}

func (u *ResumeTaskUseCase) Execute(arg interface{}) error {
	task, ok := arg.(task.Task)
	if !ok {
		return errors.New("😕 [InvalidArgumentError] arg must be `task.Task`")
	}
	logrus.Infof("🐛ResumeTaskUseCase#Execute task:%v", task)
	if err := task.Resume(time.Now()); err != nil {
		return err
	}
	u.taskRepo.Save(task)
	return nil
}
