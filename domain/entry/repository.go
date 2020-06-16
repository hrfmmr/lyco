package entry

type Repository interface {
	Add(e Entry)
	GetLatest() (Entry, error)
	GetAll() ([]Entry, error)
	Update(e Entry)
}
