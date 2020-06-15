package entry

import "errors"

type (
	Name interface {
		Value() string
	}
	entryName struct {
		value string
	}
)

func NewName(name string) (Name, error) {
	if name == "" {
		return nil, errors.New("😕 entry.Name must not be empty string")
	}
	return &entryName{name}, nil
}

func (n *entryName) Value() string {
	return n.value
}
