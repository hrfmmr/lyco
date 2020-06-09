package event

//go:generate stringer -type=EventType
type EventType int

const (
	EventTypeAny EventType = 1 << (iota + 1)
	EventTypeTaskStarted
	EventTypeTaskPaused
	EventTypeTaskResumed
	EventTypeTimerTicked
	EventTypeTimerFinished
)

type Event interface {
	Type() EventType
}
