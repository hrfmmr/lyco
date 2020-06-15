package entry

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type (
	ID interface {
		Value() string
	}

	entryid struct {
		value string
	}
)

func NewID() ID {
	id := uuid.New()
	return &entryid{id.String()}
}

func NewIDFromString(uuidStr string) (ID, error) {
	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("❗ uuidStr:%s must be uuid format", uuidStr))
	}
	return &entryid{id.String()}, nil
}

func (e *entryid) Value() string {
	return e.value
}
