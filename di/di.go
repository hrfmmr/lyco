//+build wireinject
//go:generate wire

package di

import (
	"github.com/google/wire"
	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/infra"
)

var (
	taskRepository = infra.NewTaskRepository()
)

func ProvideTaskRepository() task.Repository {
	return taskRepository
}

func provideTaskStore() store.TaskStore {
	panic(
		wire.Build(
			store.NewTaskStore,
			ProvideTaskRepository,
		),
	)
}

func provideStoreGroup() store.StoreGroup {
	panic(
		wire.Build(
			store.NewStoreGroup,
			provideTaskStore,
		),
	)
}

func InitAppContext() application.AppContext {
	panic(
		wire.Build(
			application.NewAppContext,
			provideStoreGroup,
		),
	)
}

func InitStartTaskUseCase() *usecase.StartTaskUseCase {
	panic(
		wire.Build(
			usecase.NewStartTaskUseCase,
			ProvideTaskRepository,
		),
	)
}

func InitPauseTaskUseCase() *usecase.PauseTaskUseCase {
	panic(
		wire.Build(
			usecase.NewPauseTaskUseCase,
			ProvideTaskRepository,
		),
	)
}

func InitResumeTaskUseCase() *usecase.ResumeTaskUseCase {
	panic(
		wire.Build(
			usecase.NewResumeTaskUseCase,
			ProvideTaskRepository,
		),
	)
}
