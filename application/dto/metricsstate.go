package dto

import (
	"time"

	"github.com/elliotchance/orderedmap"

	"github.com/hrfmmr/lyco/domain/entry"
)

type (
	MetricsEntry interface {
		Name() string
		Elapsed() time.Duration
		PomsCount() uint
	}

	MetricsState interface {
		Entries() []MetricsEntry
		TotalElapsed() time.Duration
		TotalPomsCount() uint
	}

	metricsEntry struct {
		entries  []entry.Entry
		duration int64
	}

	metricsState struct {
		metricsEntries []MetricsEntry
		duration       int64
	}
)

func NewMetricsEntry(entries []entry.Entry, duration int64) MetricsEntry {
	return &metricsEntry{entries, duration}
}

func (e *metricsEntry) Name() string {
	if len(e.entries) == 0 {
		return ""
	}
	return e.entries[0].Name().Value()
}

func (e *metricsEntry) Elapsed() time.Duration {
	var elapsed int64
	for _, v := range e.entries {
		elapsed += v.Elapsed().Value()
	}
	return time.Duration(elapsed)
}

func (e *metricsEntry) PomsCount() uint {
	return uint(int64(e.Elapsed()) / e.duration)
}

func NewInitialMetricsState() MetricsState {
	return &metricsState{}
}

func NewMetricsState(entries []entry.Entry, duration int64) MetricsState {
	return &metricsState{
		EntriesToMetricsModel(entries, duration),
		duration,
	}
}

func (s *metricsState) Entries() []MetricsEntry {
	return s.metricsEntries
}

func (s *metricsState) TotalElapsed() time.Duration {
	var elapsed int64
	for _, e := range s.metricsEntries {
		elapsed += int64(e.Elapsed())
	}
	return time.Duration(elapsed)
}

func (s *metricsState) TotalPomsCount() uint {
	return uint(int64(s.TotalElapsed()) / s.duration)
}

func EntriesToMetricsModel(entries []entry.Entry, duration int64) []MetricsEntry {
	m := orderedmap.NewOrderedMap()
	for _, e := range entries {
		k := e.Name().Value()
		if v, ok := m.Get(k); ok {
			s := v.([]entry.Entry)
			m.Set(k, append(s, e))
		} else {
			m.Set(k, []entry.Entry{e})
		}
	}
	metricsentries := make([]MetricsEntry, m.Len())
	for i, k := range m.Keys() {
		v, _ := m.Get(k)
		s := v.([]entry.Entry)
		metricsentries[i] = NewMetricsEntry(s, duration)
	}
	return metricsentries
}
