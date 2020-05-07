package task

import "errors"

type (
	Name interface {
		Value() string
	}
	taskName struct {
		value string
	}
)

func NewName(name string) (Name, error) {
	if name == "" {
		return nil, errors.New("😕 Task name must not be empty string")
	}
	return &taskName{name}, nil
}

func (n *taskName) Value() string {
	return n.value
}
