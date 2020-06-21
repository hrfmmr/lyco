package store

import (
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/config"
)

type MetricsStore interface {
	Store
	GetState() dto.MetricsState
	SetState(state dto.MetricsState)
}

type metricsStore struct {
	onChangeCh chan Store
	state      dto.MetricsState
}

func NewMetricsStore(cfg *config.Config) MetricsStore {
	return &metricsStore{
		make(chan Store, 1),
		dto.NewInitialMetricsState(cfg),
	}
}

func (s *metricsStore) RecvPayload(p usecase.Payload) {
}

func (s *metricsStore) OnChange() <-chan Store {
	return s.onChangeCh
}

func (s *metricsStore) GetState() dto.MetricsState {
	return s.state
}

func (s *metricsStore) SetState(state dto.MetricsState) {
	s.state = state
	s.onChangeCh <- s
}
