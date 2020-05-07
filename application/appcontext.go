package application

import (
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/sirupsen/logrus"
)

type AppContext interface {
	UseCase(usecase usecase.UseCase) usecase.UseCaseExecutor
	OnChange() <-chan store.StoreGroup
}

type appContext struct {
	storeGroup store.StoreGroup
}

func NewAppContext(storeGroup store.StoreGroup) AppContext {
	return &appContext{storeGroup}
}

func (c *appContext) UseCase(useCase usecase.UseCase) usecase.UseCaseExecutor {
	ex := usecase.NewUseCaseExecutor(useCase)
	go func(ex usecase.UseCaseExecutor) {
	Loop:
		for {
			select {
			case p := <-ex.OnWillExecute():
				logrus.Infof("👉 usecase:%v payload:%v", ex.UseCase(), p)
			case p := <-ex.OnDidExecute():
				logrus.Infof("✔ usecase:%v payload:%v", ex.UseCase(), p)
				c.storeGroup.Commit(p, usecase.NewPayLoadMeta(ex.UseCase()))
				break Loop
			}
		}
	}(ex)
	return ex
}

func (c *appContext) OnChange() <-chan store.StoreGroup {
	return c.storeGroup.OnChange()
}
