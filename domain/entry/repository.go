package entry

type Repository interface {
	Add(entry Entry)
	GetLatest() (Entry, error)
	GetAll() ([]Entry, error)
	Update(id ID, entry Entry)
}
