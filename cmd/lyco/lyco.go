package main

import (
	"sync"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/application/lifecycle"
	"github.com/hrfmmr/lyco/di"
	"github.com/hrfmmr/lyco/domain/breaks"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/ui"
	"github.com/hrfmmr/lyco/utils/notifier"
	"github.com/hrfmmr/lyco/utils/timer"

	"github.com/sirupsen/logrus"
)

const (
	longBreaksPerPoms = 4
)

var (
	wg                 sync.WaitGroup
	app                *gowid.App
	flow               gowid.RenderFlow
	err                error
	finCh              = make(chan struct{}, 1)
	startTaskUseCase   = di.InitStartTaskUseCase()
	pauseTaskUseCase   = di.InitPauseTaskUseCase()
	resumeTaskUseCase  = di.InitResumeTaskUseCase()
	stopTaskUseCase    = di.InitStopTaskUseCase()
	switchTaskUseCase  = di.InitSwitchTaskUseCase()
	abortBreaksUseCase = di.InitAbortBreaksUseCase()
	taskRepository     = di.ProvideTaskRepository()
)

func init() {
	event.DefaultPublisher.Subscribe(
		lifecycle.NewLifecycleEventHub(),
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
		appctx := di.InitAppContext()
		tasktimer := timer.NewTaskTimer()
		breakstimer := timer.NewBreaksTimer()
		fincount := 0
	Loop:
		for {
			select {
			case <-finCh:
				logrus.Infof("🔴 #main case <-finChan")
				break Loop
			case sg := <-appctx.OnChange():
				logrus.Infof("♻ #main case <-appctx.OnChange")
				task := sg.GetTask().State()
				ui.UpdateTask(app, task)
			case s := <-ui.OnStartTask():
				if t := taskRepository.GetCurrent(); t != nil && !t.CanStart() {
					continue
				}
				logrus.Infof("🚀 ui.OnStartTask::taskName=%v", s)
				taskName, err := task.NewName(s)
				if err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				if err := appctx.UseCase(startTaskUseCase).Execute(taskName); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				t := taskRepository.GetCurrent()
				tasktimer.Stop()
				tasktimer = timer.NewTaskTimer()
				tasktimer.Start(t)
			case <-ui.OnPauseTask():
				logrus.Info("|| <-ui.OnPauseTask()")
				t := taskRepository.GetCurrent()
				if !t.CanPause() {
					continue
				}
				if err := appctx.UseCase(pauseTaskUseCase).Execute(t); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				tasktimer.Stop()
			case <-ui.OnResumeTask():
				logrus.Info("▶ <-ui.OnResumeTask()")
				t := taskRepository.GetCurrent()
				if !t.CanResume() {
					continue
				}
				if err := appctx.UseCase(resumeTaskUseCase).Execute(t); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				tasktimer = timer.NewTaskTimer()
				tasktimer.Start(t)
			case <-ui.OnStopTask():
				logrus.Info("🔴 <-ui.OnStopTask()")
				t := taskRepository.GetCurrent()
				if !t.CanAbort() {
					continue
				}
				if err := appctx.UseCase(stopTaskUseCase).Execute(t); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				tasktimer.Stop()
			case s := <-ui.OnSwitchTask():
				logrus.Info("🔄 <-ui.OnSwitchTask()")
				taskName, err := task.NewName(s)
				if err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				if err := appctx.UseCase(switchTaskUseCase).Execute(taskName); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				switched := taskRepository.GetCurrent()
				tasktimer.Stop()
				tasktimer = timer.NewTaskTimer()
				tasktimer.Start(switched)
			case task := <-tasktimer.Ticker():
				logrus.Infof("⏰ #main case task := <-tasktimer.Ticker()")
				ui.UpdateTask(app, dto.ConvertTaskToDTO(task))
			case <-tasktimer.OnFinished():
				logrus.Infof("✔ #main case <-tasktimer.OnFinished()")
				fincount++
				b := breaks.ShortDefault()
				if fincount > 0 && fincount%longBreaksPerPoms == 0 {
					b = breaks.LongDefault()
				}
				if err := b.Start(); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				breakstimer = timer.NewBreaksTimer()
				breakstimer.Start(b)
				notifier.NotifyForBreaksStart(notifier.New(), b)
			case b := <-breakstimer.Ticker():
				logrus.Infof("⏰ #main case b := <-breakstimer.Ticker()")
				ui.UpdateBreaks(app, dto.ConvertBreaksToDTO(b))
			case <-breakstimer.OnFinished():
				logrus.Infof("✔ #main case <-breakstimer.OnFinished()")
				latest := taskRepository.GetCurrent()
				if err := appctx.UseCase(startTaskUseCase).Execute(latest.Name()); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				t := taskRepository.GetCurrent()
				tasktimer = timer.NewTaskTimer()
				tasktimer.Start(t)
				notifier.NotifyForBreaksEnd(notifier.New(), t)
			case <-ui.OnAbortBreaks():
				if b := breakstimer.Breaks(); b != nil {
					if b.IsStopped() {
						continue
					}
					b.Stop()
				}
				logrus.Infof("🔴 #main case <-ui.OnAbortBreaks()")
				breakstimer.Stop()
				latest := taskRepository.GetCurrent()
				if err := appctx.UseCase(abortBreaksUseCase).Execute(latest.Name()); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				t := taskRepository.GetCurrent()
				ui.UpdateTask(app, dto.ConvertTaskToDTO(t))
			}
		}
	}(app)
	ui.StartTask(app)
	app.MainLoop(gowid.UnhandledInputFunc(ui.UnhandledInput))
	finCh <- struct{}{}
	wg.Wait()
}
