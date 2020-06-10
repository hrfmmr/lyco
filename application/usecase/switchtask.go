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
	SwitchTaskPayload struct {
		taskName     string
		taskDuration time.Duration
	}

	SwitchTaskUseCase struct {
		pomodorotimer  timer.Timer
		taskService    task.TaskService
		taskRepository task.Repository
	}
)

func NewSwitchTaskPayload(name string, duration time.Duration) *SwitchTaskPayload {
	return &SwitchTaskPayload{
		name,
		duration,
	}
}

func (p *SwitchTaskPayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string `json:"name"`
		Duration string `json:"duration"`
	}{
		p.taskName,
		p.taskDuration.String(),
	})
}

func (p *SwitchTaskPayload) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

func NewSwitchTaskUseCase(
	pomodorotimer timer.Timer,
	taskService task.TaskService,
	taskRepository task.Repository,
) *SwitchTaskUseCase {
	return &SwitchTaskUseCase{
		pomodorotimer,
		taskService,
		taskRepository,
	}
}

func (u *SwitchTaskUseCase) Execute(arg interface{}) error {
	p, ok := arg.(*SwitchTaskPayload)
	if !ok {
		return errors.New(fmt.Sprintf("😕 arg:%T must be `*SwitchTaskPayload`", arg))
	}
	n, err := task.NewName(p.taskName)
	if err != nil {
		return err
	}
	d, err := task.NewDuration(int64(p.taskDuration))
	if err != nil {
		return err
	}
	t, err := u.taskService.SwitchTask(n, d)
	if err != nil {
		return err
	}
	u.taskRepository.Save(t)
	u.pomodorotimer.Start(t.Duration(), t.Elapsed())
	return nil
}
