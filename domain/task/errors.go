package task

import "fmt"

type InvalidStatusTransition struct {
	msg string
}

func NewInvalidStatusTransition(msg string) *InvalidStatusTransition {
	return &InvalidStatusTransition{msg}
}

func (e *InvalidStatusTransition) Error() string {
	return fmt.Sprintf("[InvalidStatusTransition] %s", e.msg)
}
