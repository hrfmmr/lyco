package main

import (
	"sync"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/hrfmmr/lyco/application/lifecycle"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/di"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/ui"

	"github.com/sirupsen/logrus"
)

var (
	wg                          sync.WaitGroup
	app                         *gowid.App
	flow                        gowid.RenderFlow
	err                         error
	finCh                       = make(chan struct{}, 1)
	appContext                  = di.ProvideAppContext()
	startTaskUseCase            = di.InitStartTaskUseCase()
	pauseTaskUseCase            = di.InitPauseTaskUseCase()
	resumeTaskUseCase           = di.InitResumeTaskUseCase()
	stopTaskUseCase             = di.InitStopTaskUseCase()
	switchTaskUseCase           = di.InitSwitchTaskUseCase()
	abortBreaksUseCase          = di.InitAbortBreaksUseCase()
	taskStartedEventProcessor   = di.InitTaskStartedEventProcessor()
	timertickedEventProcessor   = di.InitTimerTickedEventProcessor()
	timerfinishedEventProcessor = di.InitTimerFinishedEventProcessor()
)

func init() {
	event.DefaultPublisher.Subscribe(
		lifecycle.NewLifecycleEventHub(),
		taskStartedEventProcessor,
		timertickedEventProcessor,
		timerfinishedEventProcessor,
	)
}

func main() {
	f := examples.RedirectLogger("lyco.log")
	defer f.Close()

	app, err = ui.Build()
	examples.ExitOnErr(err)
	wg.Add(1)
	go func(app gowid.IApp) {
		defer wg.Done()
	Loop:
		for {
			select {
			case <-finCh:
				logrus.Infof("🔴 #main case <-finChan")
				break Loop
			case sg := <-appContext.OnChange():
				logrus.Infof("♻ #main case <-appctx.OnChange")
				ui.UpdatePomodoro(app, sg.TaskStore().GetState())
				ui.UpdateMetrics(app, sg.MetricsStore().GetState())
			case s := <-ui.OnStartTask():
				p := usecase.NewStartTaskPayload(s, task.DefaultDuration)
				if err := appContext.UseCase(startTaskUseCase).Execute(p); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			case <-ui.OnPauseTask():
				logrus.Info("|| <-ui.OnPauseTask()")
				if err := appContext.UseCase(pauseTaskUseCase).Execute(nil); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			case <-ui.OnResumeTask():
				logrus.Info("▶ <-ui.OnResumeTask()")
				if err := appContext.UseCase(resumeTaskUseCase).Execute(nil); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			case <-ui.OnStopTask():
				logrus.Info("🔴 <-ui.OnStopTask()")
				if err := appContext.UseCase(stopTaskUseCase).Execute(nil); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			case s := <-ui.OnSwitchTask():
				logrus.Info("🔄 <-ui.OnSwitchTask()")
				p := usecase.NewSwitchTaskPayload(s, task.DefaultDuration)
				if err := appContext.UseCase(switchTaskUseCase).Execute(p); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			case <-ui.OnAbortBreaks():
				logrus.Infof("🔴 #main case <-ui.OnAbortBreaks()")
				if err := appContext.UseCase(abortBreaksUseCase).Execute(nil); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
			}
		}
	}(app)
	ui.StartTask(app)
	app.MainLoop(gowid.UnhandledInputFunc(ui.UnhandledInput))
	finCh <- struct{}{}
	wg.Wait()
}
