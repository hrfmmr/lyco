package db

import (
	"github.com/hrfmmr/lyco/domain/task"
	log "github.com/sirupsen/logrus"
)

type (
	taskRecord struct {
		name      string
		duration  int64
		startedAt *int64
		elapsed   int64
		status    string
	}

	taskRepository struct {
		t *taskRecord
	}
)

func NewTaskRecord(name string, duration int64, startedAt *int64, elapsed int64, status string) *taskRecord {
	return &taskRecord{name, duration, startedAt, elapsed, status}
}

func NewTaskRepository() task.Repository {
	return &taskRepository{}
}

func (r *taskRepository) GetCurrent() task.Task {
	if r.t == nil {
		return nil
	}
	t, err := taskRecordToModel(r.t)
	if err != nil {
		log.Error(err)
		return nil
	}
	return t
}

func (r *taskRepository) Save(t task.Task) {
	r.t = taskModelToRecord(t)
}

func taskRecordToModel(r *taskRecord) (task.Task, error) {
	name, err := task.NewName(r.name)
	if err != nil {
		return nil, err
	}
	duration, err := task.NewDuration(r.duration)
	if err != nil {
		return nil, err
	}
	var startedAt task.StartedAt
	if r.startedAt != nil {
		startedAt, err = task.NewStartedAt(*r.startedAt)
		if err != nil {
			return nil, err
		}
	}
	elapsed, err := task.NewElapsed(r.elapsed)
	if err != nil {
		return nil, err
	}
	status, err := task.NewStatusFromString(r.status)
	if err != nil {
		return nil, err
	}
	return task.NewTaskWithValues(name, duration, startedAt, elapsed, status), nil
}

func taskModelToRecord(t task.Task) *taskRecord {
	var startedAt *int64
	if t.StartedAt() != nil {
		v := t.StartedAt().Value()
		startedAt = &v
	}
	return NewTaskRecord(
		t.Name().Value(),
		t.Duration().Value(),
		startedAt,
		t.Elapsed().Value(),
		string(t.Status().Value()),
	)
}
