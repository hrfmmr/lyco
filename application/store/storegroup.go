package store

import (
	"github.com/hrfmmr/lyco/application/usecase"
)

type (
	StoreGroup interface {
		TaskStore() TaskStore
		MetricsStore() MetricsStore
		OnChange() <-chan StoreGroup
		Commit(p usecase.Payload, meta usecase.PayloadMeta)
	}

	storegroup struct {
		taskStore    TaskStore
		metricsStore MetricsStore
		onChange     chan StoreGroup
	}
)

func NewStoreGroup(taskStore TaskStore, metricsStore MetricsStore) StoreGroup {
	sg := &storegroup{
		taskStore,
		metricsStore,
		make(chan StoreGroup),
	}
	sg.registerStores()
	return sg
}

func (g *storegroup) TaskStore() TaskStore {
	return g.taskStore
}

func (g *storegroup) MetricsStore() MetricsStore {
	return g.metricsStore
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
		g.metricsStore,
	}
}
