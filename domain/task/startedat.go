package task

type (
	StartedAt interface {
		Value() int64
	}
	taskStartedAt struct {
		value int64
	}
)

func NewStartedAt(at int64) StartedAt {
	return &taskStartedAt{at}
}

func (s *taskStartedAt) Value() int64 {
	return s.value
}
