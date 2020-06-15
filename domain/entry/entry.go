package entry

import (
	"encoding/json"
	"time"
)

type (
	Entry interface {
		// props
		ID() ID
		Name() Name
		StartedAt() StartedAt
		Elapsed() Elapsed
		// behaviors
		IncrementElapsed(elapsed Elapsed)
	}

	entry struct {
		id        ID
		name      Name
		startedAt StartedAt
		elapsed   Elapsed
	}
)

func NewEntry(name Name, startedAt StartedAt) Entry {
	id := NewID()
	elapsed, _ := NewElapsed(0)
	return &entry{id, name, startedAt, elapsed}
}

func NewEntryWithValues(id ID, name Name, startedAt StartedAt, elapsed Elapsed) Entry {
	return &entry{id, name, startedAt, elapsed}
}

func (e *entry) ID() ID {
	return e.id
}

func (e *entry) Name() Name {
	return e.name
}

func (e *entry) StartedAt() StartedAt {
	return e.startedAt
}

func (e *entry) Elapsed() Elapsed {
	return e.elapsed
}

func (e *entry) IncrementElapsed(elapsed Elapsed) {
	e.elapsed.Add(elapsed)
}

func (e *entry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string `json:"name"`
		StartedAt string `json:"started_at"`
		Elapsed   string `json:"elapsed"`
	}{
		e.name.Value(),
		time.Unix(0, e.startedAt.Value()).String(),
		time.Duration(e.elapsed.Value()).String(),
	})
}

func (e *entry) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}
