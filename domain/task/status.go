package task

import (
	"errors"
	"fmt"
)

type StatusValue string

const (
	TaskStatusNone     StatusValue = "none"
	TaskStatusRunning  StatusValue = "running"
	TaskStatusPaused   StatusValue = "paused"
	TaskStatusStopped  StatusValue = "aborted"
	TaskStatusFinished StatusValue = "finished"
)

type Status interface {
	Value() StatusValue
	Update(to Status) error
}

type status struct {
	value StatusValue
}

func NewStatus(s StatusValue) Status {
	return &status{s}
}

func NewStatusFromString(s string) (Status, error) {
	switch s {
	case string(TaskStatusNone):
		return NewStatus(TaskStatusNone), nil
	case string(TaskStatusRunning):
		return NewStatus(TaskStatusRunning), nil
	case string(TaskStatusPaused):
		return NewStatus(TaskStatusPaused), nil
	case string(TaskStatusStopped):
		return NewStatus(TaskStatusStopped), nil
	case string(TaskStatusFinished):
		return NewStatus(TaskStatusFinished), nil
	}
	return nil, errors.New(fmt.Sprintf("❗got invalid status:%s", s))
}

func (s *status) Value() StatusValue {
	return s.value
}

// [TaskStatus state diagram](http://www.plantuml.com/plantuml/png/VP5DIi0m48NtESLGDohq0YvA1N4fKfUbYsZ6DjXEIduM7i1RU3fFuYJDK0IrYpAPRoRlFTA7g7rCsweMQn1ms-Cx60mltkxHEbBC8qBpu0WRq0682saEYSZINFh-g0KzwXJG5BANKi2z9HkMYbxhGU3ji_EnCdtKIetN4xHjToZdNpw97jp0KqvmUQMaquuNiqUaYQT4WFSYWYDOBRqfuE_E4NvzVYWa0ncUwrrrp5UN57nrRVoG7J2axOHPrgDXKN7ECoNFQUtV6R7wHpfP-9tsI8OVyHi0)
func (s *status) Update(to Status) error {
	for _, v := range nextStatuses(s.value) {
		if to.Value() == v {
			s.value = to.Value()
			return nil
		}
	}
	return NewInvalidStatusTransition(fmt.Sprintf("not allowed state transition from %s to %s", s.Value(), to.Value()))
}

func nextStatuses(s StatusValue) []StatusValue {
	switch s {
	case TaskStatusNone:
		return []StatusValue{
			TaskStatusRunning,
		}
	case TaskStatusRunning:
		return []StatusValue{
			TaskStatusPaused,
			TaskStatusFinished,
			TaskStatusStopped,
		}
	case TaskStatusPaused:
		return []StatusValue{
			TaskStatusRunning,
			TaskStatusStopped,
		}
	default:
		return nil
	}
}
