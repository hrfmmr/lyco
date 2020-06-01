package event

type Event interface {
	Type() EventType
}
