package task

import (
	"fmt"
)

type StatusValue string

const (
	TaskStatusNone     StatusValue = "none"
	TaskStatusRunning  StatusValue = "running"
	TaskStatusPaused   StatusValue = "paused"
	TaskStatusAborted  StatusValue = "aborted"
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
			TaskStatusAborted,
		}
	case TaskStatusPaused:
		return []StatusValue{
			TaskStatusRunning,
			TaskStatusAborted,
		}
	default:
		return nil
	}
}
