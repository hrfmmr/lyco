// Code generated by "stringer -type=AvailableTaskAction"; DO NOT EDIT.

package dto

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AvailableTaskActionStart-0]
	_ = x[AvailableTaskActionPause-1]
	_ = x[AvailableTaskActionResume-2]
	_ = x[AvailableTaskActionStop-3]
	_ = x[AvailableTaskActionSwitch-4]
	_ = x[AvailableTaskActionAbortBreaks-5]
}

const _AvailableTaskAction_name = "AvailableTaskActionStartAvailableTaskActionPauseAvailableTaskActionResumeAvailableTaskActionStopAvailableTaskActionSwitchAvailableTaskActionAbortBreaks"

var _AvailableTaskAction_index = [...]uint8{0, 24, 48, 73, 96, 121, 151}

func (i AvailableTaskAction) String() string {
	if i < 0 || i >= AvailableTaskAction(len(_AvailableTaskAction_index)-1) {
		return "AvailableTaskAction(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _AvailableTaskAction_name[_AvailableTaskAction_index[i]:_AvailableTaskAction_index[i+1]]
}
