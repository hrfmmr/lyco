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
	NextStatuses() []StatusValue
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

// [TaskStatus state diagram](http://www.plantuml.com/plantuml/png/ZP9DJiD038NtFeNLRW0f5s1H5GcnHAWiemeDZ6b7D9veVXOSW5jmEaw2CyLIJ5YGLUQzB_5xrcnpSQdKuGHm39wGXh6yewVyai9OGcGGe12kxYFJ2bt6TdvYEQgrgyo13pCtdHK57bpDv6V-s0IrxmA7V3J0wu-aoCrpJCKGxgm0z5TxxBgDpQMlpJ6Py1eVfyeNbs3rhbyV4X7lsnSc1CT261bFWTy0vRcDJ5-V7q3iSJ9__sfZUih8jW4TgV8VqNQKwx04-tKYBy5iT-5bByig_-QrODkIx14ihkgMs4ytv1i0)
func (s *status) Update(to Status) error {
	for _, v := range nextStatuses(s.value) {
		if to.Value() == v {
			s.value = to.Value()
			return nil
		}
	}
	return NewInvalidStatusTransition(fmt.Sprintf("not allowed state transition from %s to %s", s.Value(), to.Value()))
}

func (s *status) NextStatuses() []StatusValue {
	return nextStatuses(s.value)
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
