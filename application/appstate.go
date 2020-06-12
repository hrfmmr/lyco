package application

import "github.com/hrfmmr/lyco/domain/breaks"

type (
	AppState interface {
		CurrentBreaks() breaks.Breaks
	}

	appstate struct {
		b breaks.Breaks
	}
)

func NewAppState() AppState {
	return &appstate{}
}

func (s *appstate) CurrentBreaks() breaks.Breaks {
	return s.b
}

func (s *appstate) SetBreaks(b breaks.Breaks) {
	s.b = b
}
