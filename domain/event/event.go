package event

//go:generate stringer -type=EventType
type EventType int

const (
	EventTypeAny EventType = 1 << (iota + 1)
	EventTypeTaskStarted
)

type Event interface {
	Type() EventType
}
