package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gcla/gowid"
	"github.com/hrfmmr/lyco/application/lifecycle"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/di"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/ui"
	"github.com/shibukawa/configdir"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
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
	os.Exit(cmain())
}

func cmain() int {
	stdConf := configdir.New("", "lyco")
	dirs := stdConf.QueryFolders(configdir.Cache)
	if err := dirs[0].CreateParentDir("dummy"); err != nil {
		fmt.Printf("Warning: could not create cache dir: %v\n", err)
	}
	cachedir := dirs[0].Path
	logfile := filepath.Join(cachedir, "lyco.log")
	logfd, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create log file %s: %v\n", logfile, err)
		return 1
	}
	// Don't close it - just let the descriptor be closed at exit. log is used
	// in many places, some outside of this main function, and closing results in
	// an error often on freebsd.
	//defer logfd.Close()
	log.SetOutput(logfd)

	app, err = ui.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed initializing app err:%v\n", err)
		return 1
	}
	var g errgroup.Group
	g.Go(func() error {
		for {
			select {
			case <-finCh:
				return nil
			case sg := <-appContext.OnChange():
				ui.UpdatePomodoro(app, sg.TaskStore().GetState())
				ui.UpdateMetrics(app, sg.MetricsStore().GetState())
			case s := <-ui.OnStartTask():
				p := usecase.NewStartTaskPayload(s, task.DefaultDuration)
				if err := appContext.UseCase(startTaskUseCase).Execute(p); err != nil {
					return err
				}
			case <-ui.OnPauseTask():
				if err := appContext.UseCase(pauseTaskUseCase).Execute(nil); err != nil {
					return err
				}
			case <-ui.OnResumeTask():
				if err := appContext.UseCase(resumeTaskUseCase).Execute(nil); err != nil {
					return err
				}
			case <-ui.OnStopTask():
				if err := appContext.UseCase(stopTaskUseCase).Execute(nil); err != nil {
					return err
				}
			case s := <-ui.OnSwitchTask():
				p := usecase.NewSwitchTaskPayload(s, task.DefaultDuration)
				if err := appContext.UseCase(switchTaskUseCase).Execute(p); err != nil {
					return err
				}
			case <-ui.OnAbortBreaks():
				if err := appContext.UseCase(abortBreaksUseCase).Execute(nil); err != nil {
					return err
				}
			}
		}
	})
	ui.StartTask(app)
	app.MainLoop(gowid.UnhandledInputFunc(ui.UnhandledInput))
	finCh <- struct{}{}
	if err := g.Wait(); err != nil {
		log.Error(err)
		return 1
	}
	return 0
}
