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
	"github.com/hrfmmr/lyco/domain/entry"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/hrfmmr/lyco/infra/db"
)

var (
	pomodorotimer   = timer.NewTimer()
	taskRepository  = db.NewTaskRepository()
	entryRepository = db.NewEntryRepository()
	taskStore       = store.NewTaskStore(
		taskRepository,
	)
	metricsStore = store.NewMetricsStore()
	storeGroup   = store.NewStoreGroup(
		taskStore,
		metricsStore,
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

func provideMetricsStore() store.MetricsStore {
	return metricsStore
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

func provideEntryRepository() entry.Repository {
	return entryRepository
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

func InitTaskStartedEventProcessor() *eventprocessor.TaskStartedEventProcessor {
	panic(
		wire.Build(
			eventprocessor.NewTaskStartedEventProcessor,
			provideEntryRepository,
		),
	)
}

func InitTimerTickedEventProcessor() *eventprocessor.TimerTickedEventProcessor {
	panic(
		wire.Build(
			eventprocessor.NewTimerTickedEventProcessor,
			ProvideAppContext,
			provideMetricsStore,
			ProvideTaskRepository,
			provideEntryRepository,
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
