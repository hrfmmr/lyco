package store

import (
	"github.com/hrfmmr/lyco/application/usecase"
)

type Store interface {
	RecvPayload(p usecase.Payload)
	OnChange() <-chan Store
}
