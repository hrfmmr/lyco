// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/hrfmmr/lyco/application"
	"github.com/hrfmmr/lyco/application/appstate"
	"github.com/hrfmmr/lyco/application/eventprocessor"
	"github.com/hrfmmr/lyco/application/store"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/config"
	"github.com/hrfmmr/lyco/domain/entry"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/domain/timer"
	"github.com/hrfmmr/lyco/infra/db"
)

// Injectors from di.go:

func provideTaskService() task.TaskService {
	repository := ProvideTaskRepository()
	taskService := task.NewTaskService(repository)
	return taskService
}

func InitStartTaskUseCase() *usecase.StartTaskUseCase {
	timer := provideTimer()
	repository := ProvideTaskRepository()
	startTaskUseCase := usecase.NewStartTaskUseCase(timer, repository)
	return startTaskUseCase
}

func InitPauseTaskUseCase() *usecase.PauseTaskUseCase {
	timer := provideTimer()
	repository := ProvideTaskRepository()
	pauseTaskUseCase := usecase.NewPauseTaskUseCase(timer, repository)
	return pauseTaskUseCase
}

func InitResumeTaskUseCase() *usecase.ResumeTaskUseCase {
	timer := provideTimer()
	repository := ProvideTaskRepository()
	resumeTaskUseCase := usecase.NewResumeTaskUseCase(timer, repository)
	return resumeTaskUseCase
}

func InitStopTaskUseCase() *usecase.StopTaskUseCase {
	timer := provideTimer()
	repository := ProvideTaskRepository()
	stopTaskUseCase := usecase.NewStopTaskUseCase(timer, repository)
	return stopTaskUseCase
}

func InitSwitchTaskUseCase() *usecase.SwitchTaskUseCase {
	timer := provideTimer()
	taskService := provideTaskService()
	repository := ProvideTaskRepository()
	switchTaskUseCase := usecase.NewSwitchTaskUseCase(timer, taskService, repository)
	return switchTaskUseCase
}

func InitAbortBreaksUseCase() *usecase.AbortBreaksUseCase {
	appstateAppState := provideAppState()
	repository := ProvideTaskRepository()
	timer := provideTimer()
	abortBreaksUseCase := usecase.NewAbortBreaksUseCase(appstateAppState, repository, timer)
	return abortBreaksUseCase
}

func InitTaskStartedEventProcessor() *eventprocessor.TaskStartedEventProcessor {
	repository := provideEntryRepository()
	taskStartedEventProcessor := eventprocessor.NewTaskStartedEventProcessor(repository)
	return taskStartedEventProcessor
}

func InitTimerTickedEventProcessor() *eventprocessor.TimerTickedEventProcessor {
	applicationAppContext := ProvideAppContext()
	storeMetricsStore := provideMetricsStore()
	repository := ProvideTaskRepository()
	repository2 := provideEntryRepository()
	timerTickedEventProcessor := eventprocessor.NewTimerTickedEventProcessor(applicationAppContext, storeMetricsStore, repository, repository2)
	return timerTickedEventProcessor
}

func InitTimerFinishedEventProcessor() *eventprocessor.TimerFinishedEventProcessor {
	config := ProvideConfig()
	applicationAppContext := ProvideAppContext()
	repository := ProvideTaskRepository()
	timer := provideTimer()
	timerFinishedEventProcessor := eventprocessor.NewTimerFinishedEventProcessor(config, applicationAppContext, repository, timer)
	return timerFinishedEventProcessor
}

// di.go:

var (
	cfg             = config.NewConfig()
	pomodorotimer   = timer.NewTimer()
	taskRepository  = db.NewTaskRepository()
	entryRepository = db.NewEntryRepository()
	taskStore       = store.NewTaskStore(
		taskRepository,
	)
	metricsStore = store.NewMetricsStore(cfg)
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

func ProvideConfig() *config.Config {
	return cfg
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
