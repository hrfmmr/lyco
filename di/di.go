//+build wireinject
//go:generate wire

package di

import (
	"github.com/google/wire"
	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/appstate"
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
	storeGroup = store.NewStoreGroup(
		taskStore,
	)
	appState   = appstate.NewAppState()
	appContext = application.NewAppContext(
		appState,
		storeGroup,
	)
	eventPublisher = event.NewPublisher()
)

func provideTimer() timer.Timer {
	return pomodorotimer
}

func provideTaskStore() store.TaskStore {
	return taskStore
}

func provideAppState() appstate.AppState {
	return appState
}

func provideTaskService() task.TaskService {
	panic(
		wire.Build(
			task.NewTaskService,
			ProvideTaskRepository,
		),
	)
}

func ProvideAppContext() application.AppContext {
	return appContext
}

func ProvideTaskRepository() task.Repository {
	return taskRepository
}

func InitStartTaskUseCase() *usecase.StartTaskUseCase {
	panic(
		wire.Build(
			usecase.NewStartTaskUseCase,
			provideTimer,
			ProvideTaskRepository,
		),
	)
}

func InitPauseTaskUseCase() *usecase.PauseTaskUseCase {
	panic(
		wire.Build(
			usecase.NewPauseTaskUseCase,
			provideTimer,
			ProvideTaskRepository,
		),
	)
}

func InitResumeTaskUseCase() *usecase.ResumeTaskUseCase {
	panic(
		wire.Build(
			usecase.NewResumeTaskUseCase,
			provideTimer,
			ProvideTaskRepository,
		),
	)
}

func InitStopTaskUseCase() *usecase.StopTaskUseCase {
	panic(
		wire.Build(
			usecase.NewStopTaskUseCase,
			provideTimer,
			ProvideTaskRepository,
		),
	)
}

func InitSwitchTaskUseCase() *usecase.SwitchTaskUseCase {
	panic(
		wire.Build(
			usecase.NewSwitchTaskUseCase,
			provideTimer,
			provideTaskService,
			ProvideTaskRepository,
		),
	)
}

func InitAbortBreaksUseCase() *usecase.AbortBreaksUseCase {
	panic(
		wire.Build(
			usecase.NewAbortBreaksUseCase,
			provideAppState,
			ProvideTaskRepository,
			provideTimer,
		),
	)
}

func InitTimerTickedEventProcessor() *eventprocessor.TimerTickedEventProcessor {
	panic(
		wire.Build(
			eventprocessor.NewTimerTickedEventProcessor,
			ProvideAppContext,
			ProvideTaskRepository,
		),
	)
}

func InitTimerFinishedEventProcessor() *eventprocessor.TimerFinishedEventProcessor {
	panic(
		wire.Build(
			eventprocessor.NewTimerFinishedEventProcessor,
			ProvideAppContext,
			ProvideTaskRepository,
			provideTimer,
		),
	)
}
