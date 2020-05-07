package infra

import (
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/sirupsen/logrus"
)

type taskRepository struct {
	t task.Task
}

func NewTaskRepository() task.Repository {
	return &taskRepository{}
}

func (r *taskRepository) GetCurrent() task.Task {
	return r.t
}

func (r *taskRepository) Save(t task.Task) {
	logrus.Infof("🐛taskRepository#Save t:%v", t)
	r.t = t
}
