package task_test

import (
	"testing"

	"github.com/hrfmmr/lyco/domain/task"
)

func TestStatusUpdate(t *testing.T) {
	tests := []struct {
		name     string
		from     task.Status
		to       task.Status
		expected map[string]interface{}
	}{
		{
			name: "🔨Test update status from:TaskStatusNone to:TaskStatusRunning",
			from: task.NewStatus(task.TaskStatusNone),
			to:   task.NewStatus(task.TaskStatusRunning),
			expected: map[string]interface{}{
				"error":   nil,
				"updated": task.TaskStatusRunning,
			},
		},
		{
			name: "🔨Test update status from:TaskStatusNone to:TaskStatusFinished",
			from: task.NewStatus(task.TaskStatusNone),
			to:   task.NewStatus(task.TaskStatusFinished),
			expected: map[string]interface{}{
				"error":   task.NewInvalidStatusTransition("not allowed state transition from none to finished"),
				"updated": task.TaskStatusNone,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.from.Update(test.to)
			if err != nil && err.Error() != test.expected["error"].(error).Error() {
				t.Errorf("got err:%v, expected:%v", err, test.expected["error"])
			}
			if test.from.Value() != test.expected["updated"] {
				t.Errorf("got updated:%s, expected:%s", test.from, test.expected["updated"])
			}
		})
	}
}
