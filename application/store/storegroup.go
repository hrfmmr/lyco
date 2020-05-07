package store

import (
	"github.com/hrfmmr/lyco/application/usecase"
)

type StoreGroup interface {
	GetTask() TaskState
	OnChange() <-chan StoreGroup
	Commit(p usecase.Payload, meta usecase.PayloadMeta)
}

type storeGroup struct {
	taskStore TaskStore
	onChange  chan StoreGroup
}

func NewStoreGroup(taskStore TaskStore) StoreGroup {
	return &storeGroup{
		taskStore,
		make(chan StoreGroup),
	}
}

func (g *storeGroup) GetTask() TaskState {
	return g.taskStore.GetState()
}

func (g *storeGroup) OnChange() <-chan StoreGroup {
	return g.onChange
}

func (g *storeGroup) Commit(p usecase.Payload, meta usecase.PayloadMeta) {
	for _, s := range g.stores() {
		s.RecvPayload(p)
	}
	g.onChange <- g
}

func (g *storeGroup) stores() []Store {
	return []Store{
		g.taskStore,
	}
}
