// Code generated by "stringer -type=EventType"; DO NOT EDIT.

package event

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EventTypeAny-2]
	_ = x[EventTypeTaskStarted-4]
}

const (
	_EventType_name_0 = "EventTypeAny"
	_EventType_name_1 = "EventTypeTaskStarted"
)

func (i EventType) String() string {
	switch {
	case i == 2:
		return _EventType_name_0
	case i == 4:
		return _EventType_name_1
	default:
		return "EventType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}