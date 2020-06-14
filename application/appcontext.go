package application

import (
	"github.com/hrfmmr/lyco/application/appstate"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/sirupsen/logrus"
)

type AppContext interface {
	AppState() appstate.AppState
	StoreGroup() store.StoreGroup
	UseCase(usecase usecase.UseCase) usecase.UseCaseExecutor
	OnChange() <-chan store.StoreGroup
}

type appContext struct {
	appstate   appstate.AppState
	storeGroup store.StoreGroup
}

func NewAppContext(appstate appstate.AppState, storeGroup store.StoreGroup) AppContext {
	return &appContext{
		appstate,
		storeGroup,
	}
}

func (c *appContext) AppState() appstate.AppState {
	return c.appstate
}

func (c *appContext) StoreGroup() store.StoreGroup {
	return c.storeGroup
}

func (c *appContext) UseCase(useCase usecase.UseCase) usecase.UseCaseExecutor {
	ex := usecase.NewUseCaseExecutor(useCase)
	go func(ex usecase.UseCaseExecutor) {
	Loop:
		for {
			select {
			case p := <-ex.OnWillExecute():
				logrus.Infof("👉 usecase:%T payload:%v", ex.UseCase(), p)
			case p := <-ex.OnDidExecute():
				logrus.Infof("✔ usecase:%T payload:%v", ex.UseCase(), p)
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
