package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/gcla/gowid"
	"github.com/hrfmmr/lyco/application/lifecycle"
	"github.com/hrfmmr/lyco/application/usecase"
	"github.com/hrfmmr/lyco/cli"
	"github.com/hrfmmr/lyco/di"
	"github.com/hrfmmr/lyco/domain/event"
	"github.com/hrfmmr/lyco/domain/task"
	"github.com/hrfmmr/lyco/ui"
	flags "github.com/jessevdk/go-flags"
	"github.com/shibukawa/configdir"
	"golang.org/x/sync/errgroup"

	log "github.com/sirupsen/logrus"
)

var (
	cfg                         = di.ProvideConfig()
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

func newOptsParser(opt *cli.Lyco) *flags.Parser {
	p := flags.NewParser(opt, flags.PassDoubleDash|flags.PrintErrors)
	p.Name = "lyco"
	p.Usage = "[OPTIONS]"
	return p
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
	log.SetReportCaller(true)

	var opts cli.Lyco
	parser := newOptsParser(&opts)
	_, err = parser.Parse()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Command-line error:%v\n\n", err)
		return 1
	}

	if opts.Help {
		cli.WriteHelp(parser, os.Stdout)
		return 0
	}

	if opts.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("🐛Starting debug mode")
	}

	if len(opts.PomodoroDuration) > 0 {
		d, err := time.ParseDuration(opts.PomodoroDuration)
		log.Debugf("🐛 opts.PomodoroDuration:%s", d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse relative time for pomodoro duration\n\n")
			return 1
		}
		cfg.PomodoroDuration = &d
	}

	if len(opts.ShortBreaksDuration) > 0 {
		d, err := time.ParseDuration(opts.ShortBreaksDuration)
		log.Debugf("🐛 opts.ShortBreaksDuration:%s", d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse relative time for short-breaks duration\n\n")
			return 1
		}
		cfg.ShortBreaksDuration = &d
	}

	if len(opts.LongBreaksDuration) > 0 {
		d, err := time.ParseDuration(opts.LongBreaksDuration)
		log.Debugf("🐛 opts.LongBreaksDuration:%s", d)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not parse relative time for long-breaks duration\n\n")
			return 1
		}
		cfg.LongBreaksDuration = &d
	}

	app, err := ui.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed initializing app err:%v\n", err)
		return 1
	}
	var g errgroup.Group
	finCh := make(chan struct{}, 1)
	g.Go(func() error {
		defer func() {
			if e := recover(); e != nil {
				app.Quit()
				t := time.NewTicker(time.Second)
				<-t.C
				fmt.Fprintf(os.Stderr, "Unexpected error occurred...\nSee %s for more details\n\n", logfile)
				log.Fatalf("😇 panic:  %v", string(debug.Stack()))
			}
		}()
		for {
			select {
			case <-finCh:
				return nil
			case sg := <-appContext.OnChange():
				ui.UpdatePomodoro(app, sg.TaskStore().GetState())
				ui.UpdateMetrics(app, sg.MetricsStore().GetState())
			case s := <-ui.OnStartTask():
				var d time.Duration
				if cfg.PomodoroDuration != nil {
					d = *cfg.PomodoroDuration
				} else {
					d = task.DefaultDuration
				}
				p := usecase.NewStartTaskPayload(s, d)
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
				var d time.Duration
				if cfg.PomodoroDuration != nil {
					d = *cfg.PomodoroDuration
				} else {
					d = task.DefaultDuration
				}
				p := usecase.NewSwitchTaskPayload(s, d)
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
