package db

import "github.com/hrfmmr/lyco/domain/entry"

type (
	entryRecord struct {
		id        string
		name      string
		startedAt int64
		elapsed   int64
	}

	entryRepository struct {
		entries []*entryRecord
	}
)

func NewEntryRecord(id, name string, startedAt, elapsed int64) *entryRecord {
	return &entryRecord{id, name, startedAt, elapsed}
}

func NewEntryRepository() entry.Repository {
	return &entryRepository{
		entries: []*entryRecord{},
	}
}

func (r *entryRepository) Add(e entry.Entry) {
	record := entryModelToRecord(e)
	r.entries = append(r.entries, record)
}

func (r *entryRepository) GetLatest() (entry.Entry, error) {
	if len(r.entries) == 0 {
		return nil, nil
	}
	curr := r.entries[len(r.entries)-1]
	return entryRecordToModel(curr)
}

func (r *entryRepository) GetAll() ([]entry.Entry, error) {
	entries := make([]entry.Entry, len(r.entries))
	for i, e := range r.entries {
		m, err := entryRecordToModel(e)
		if err != nil {
			return nil, err
		}
		entries[i] = m
	}
	return entries, nil
}

func (r *entryRepository) Update(e entry.Entry) {
	for i, v := range r.entries {
		if v.id == e.ID().Value() {
			r.entries[i] = entryModelToRecord(e)
			return
		}
	}
}

func entryRecordToModel(r *entryRecord) (entry.Entry, error) {
	id, err := entry.NewIDFromString(r.id)
	if err != nil {
		return nil, err
	}
	name, err := entry.NewName(r.name)
	if err != nil {
		return nil, err
	}
	startedAt, err := entry.NewStartedAt(r.startedAt)
	if err != nil {
		return nil, err
	}
	elapsed, err := entry.NewElapsed(r.elapsed)
	if err != nil {
		return nil, err
	}
	return entry.NewEntryWithValues(id, name, startedAt, elapsed), nil
}

func entryModelToRecord(e entry.Entry) *entryRecord {
	return NewEntryRecord(
		e.ID().Value(),
		e.Name().Value(),
		e.StartedAt().Value(),
		e.Elapsed().Value(),
	)
}
