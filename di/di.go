//+build wireinject
//go:generate wire

package di

import (
	"github.com/google/wire"
	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/eventprocessor"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/hrfmmr/lyco/infra"
)

var (
	pomodorotimer  = timer.NewTimer()
	taskRepository = infra.NewTaskRepository()
	taskStore      = store.NewTaskStore(
		taskRepository,
	)
	eventPublisher = event.NewPublisher()
)

func provideTimer() timer.Timer {
	return pomodorotimer
}

func ProvideTaskRepository() task.Repository {
	return taskRepository
}

func provideTaskStore() store.TaskStore {
	return taskStore
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

func InitStopTaskUseCase() *usecase.StopTaskUseCase {
	panic(
		wire.Build(
			usecase.NewStopTaskUseCase,
			ProvideTaskRepository,
		),
	)
}

func InitSwitchTaskUseCase() *usecase.SwitchTaskUseCase {
	panic(
		wire.Build(
			usecase.NewSwitchTaskUseCase,
			ProvideTaskRepository,
		),
	)
}

func InitAbortBreaksUseCase() *usecase.AbortBreaksUseCase {
	panic(
		wire.Build(
			usecase.NewAbortBreaksUseCase,
			ProvideTaskRepository,
		),
	)
}

func InitTimerStarter() *eventprocessor.TimerStarter {
	panic(
		wire.Build(
			eventprocessor.NewTimerStarter,
			provideTimer,
		),
	)
}

func InitTimerStateUpdater() *eventprocessor.TimerStateUpdater {
	panic(
		wire.Build(
			eventprocessor.NewTimerStateUpdater,
			provideTaskStore,
			ProvideTaskRepository,
		),
	)
}
