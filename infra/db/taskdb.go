package db

import (
	"github.com/hrfmmr/lyco/domain/task"
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
	r.t = t
}
