package dto_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/domain/entry"
)

func TestEntriesToMetricsModel(t *testing.T) {
	testStartedAt, err := entry.NewStartedAt(time.Now().UnixNano())
	if err != nil {
		t.Fatal(err)
	}
	testDuration := int64(time.Second)
	tests := []struct {
		name       string
		entryNames []string
		expected   map[string]interface{}
	}{
		{
			"🔨[1]Test the MetricsEntry is stored in the order of entries",
			[]string{
				"test1",
				"test2",
				"test1",
				"test2",
				"test3",
				"test2",
			},
			map[string]interface{}{
				"namesorder": []string{
					"test1",
					"test2",
					"test3",
				},
			},
		},
		{
			"🔨[2]Test the MetricsEntry is stored in the order of entries",
			[]string{
				"test5",
				"test4",
				"test4",
				"test3",
				"test5",
				"test2",
				"test3",
				"test1",
				"test2",
				"test3",
				"test5",
				"test1",
			},
			map[string]interface{}{
				"namesorder": []string{
					"test5",
					"test4",
					"test3",
					"test2",
					"test1",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entries := make([]entry.Entry, len(test.entryNames))
			for i, n := range test.entryNames {
				name, err := entry.NewName(n)
				if err != nil {
					t.Fatal(err)
				}
				e := entry.NewEntry(name, testStartedAt)
				entries[i] = e
			}
			metricsentries := dto.EntriesToMetricsModel(entries, testDuration)
			names := make([]string, len(metricsentries))
			for i, e := range metricsentries {
				names[i] = e.Name()
			}
			if !reflect.DeepEqual(names, test.expected["namesorder"]) {
				t.Errorf("❗got names:%v, expecting:%v", names, test.expected["namesorder"])
			}
		})
	}
}
