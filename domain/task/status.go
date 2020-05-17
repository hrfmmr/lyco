package task

type Status string

const (
	TaskStatusNone     Status = "none"
	TaskStatusRunning  Status = "running"
	TaskStatusPaused   Status = "paused"
	TaskStatusAborted  Status = "aborted"
	TaskStatusFinished Status = "finished"
)
