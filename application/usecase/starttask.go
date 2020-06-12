package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
)

type (
	StartTaskPayload struct {
		taskName     string
		taskDuration time.Duration
	}

	StartTaskUseCase struct {
		pomodorotimer  timer.Timer
		taskRepository task.Repository
	}
)

func NewStartTaskPayload(name string, duration time.Duration) *StartTaskPayload {
	return &StartTaskPayload{
		name,
		duration,
	}
}

func (p *StartTaskPayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string `json:"name"`
		Duration string `json:"duration"`
	}{
		p.taskName,
		p.taskDuration.String(),
	})
}

func (p *StartTaskPayload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewStartTaskUseCase(pomodorotimer timer.Timer, taskRepository task.Repository) *StartTaskUseCase {
	return &StartTaskUseCase{
		pomodorotimer,
		taskRepository,
	}
}

func (u *StartTaskUseCase) Execute(arg interface{}) error {
	if t := u.taskRepository.GetCurrent(); t != nil && !t.CanStart() {
		// ignore
		return nil
	}
	p, ok := arg.(*StartTaskPayload)
	if !ok {
		return errors.New(fmt.Sprintf("😕 arg:%T must be `*StartTaskPayload`", arg))
	}
	n, err := task.NewName(p.taskName)
	if err != nil {
		return err
	}
	d, err := task.NewDuration(int64(p.taskDuration))
	if err != nil {
		return err
	}
	t := task.NewTask(n, d)
	if err := t.Start(time.Now()); err != nil {
		return err
	}
	u.taskRepository.Save(t)
	d, err = timer.NewDuration(t.Duration().Value())
	if err != nil {
		return err
	}
	u.pomodorotimer.Start(d)
	return nil
}
