package main

import (
	"sync"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/examples"
	"github.com/hrfmmr/lyco/application/dto"
	"github.com/hrfmmr/lyco/di"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/ui"
	"github.com/hrfmmr/lyco/utils/timer"

	"github.com/sirupsen/logrus"
)

var (
	wg               sync.WaitGroup
	app              *gowid.App
	flow             gowid.RenderFlow
	err              error
	finCh            = make(chan struct{}, 1)
	startTaskUseCase = di.InitStartTaskUseCase()
	taskRepository   = di.ProvideTaskRepository()
)

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
	Loop:
		for {
			select {
			case <-finCh:
				logrus.Infof("🔴 #main case <-finChan")
				break Loop
			case sg := <-appctx.OnChange():
				logrus.Infof("♻ #main case <-appctx.OnChange")
				task := sg.GetTask().State()
				ui.Update(app, task)
			case s := <-ui.OnStartTask():
				logrus.Infof("🐛 ui.OnStartTask::taskName=%v", s)
				taskName, err := task.NewName(s)
				if err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				if err := appctx.UseCase(startTaskUseCase).Execute(taskName); err != nil {
					logrus.Fatalf("💀 %v", err)
				}
				task := taskRepository.GetCurrent()
				tasktimer.Stop()
				tasktimer = timer.NewTaskTimer()
				tasktimer.Start(task)
			case task := <-tasktimer.Ticker():
				logrus.Infof("♻ #main case task := <-tasktimer.Ticker()")
				ui.Update(app, dto.ConvertTaskToDTO(task))
			}
		}
	}(app)

	//TODO: temp
	ui.SwitchTask(app)

	app.SimpleMainLoop()
	finCh <- struct{}{}
	wg.Wait()
}
