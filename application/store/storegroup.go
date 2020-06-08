package store

import (
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/usecase"
)

type (
	StoreGroup interface {
		GetTask() dto.TaskState
		OnChange() <-chan StoreGroup
		Commit(p usecase.Payload, meta usecase.PayloadMeta)
	}

	storegroup struct {
		taskStore TaskStore
		onChange  chan StoreGroup
	}
)

func NewStoreGroup(taskStore TaskStore) StoreGroup {
	sg := &storegroup{
		taskStore,
		make(chan StoreGroup),
	}
	sg.registerStores()
	return sg
}

func (g *storegroup) GetTask() dto.TaskState {
	return g.taskStore.GetState()
}

func (g *storegroup) OnChange() <-chan StoreGroup {
	return g.onChange
}

func (g *storegroup) Commit(p usecase.Payload, meta usecase.PayloadMeta) {
	for _, s := range g.stores() {
		s.RecvPayload(p)
	}
	g.onChange <- g
}

func (g *storegroup) registerStores() {
	for _, s := range g.stores() {
		go func(store Store) {
			for {
				select {
				case <-store.OnChange():
					g.onChange <- g
				}
			}
		}(s)
	}
}

func (g *storegroup) stores() []Store {
	return []Store{
		g.taskStore,
	}
}
