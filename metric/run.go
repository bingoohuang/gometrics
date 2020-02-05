package metric

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	"github.com/bingoohuang/gometrics/rotate"
	"github.com/bingoohuang/gometrics/util"
	"github.com/sirupsen/logrus"
)

// DefaultRunner is the default runner for metric recording
var DefaultRunner *Runner // nolint

func init() { // nolint
	DefaultRunner = NewRunnerOptions(EnvOption())
	DefaultRunner.Start()
}

// Runner is a runner for metric rotate writing
type Runner struct {
	startTime time.Time
	Interval  time.Duration
	C         chan Line
	stop      chan bool
	Logfile   *rotate.File
}

// NewRunnerOptions creates a Runner
func NewRunnerOptions(o *Option) *Runner {
	f := filepath.Join(o.LogPath, "metrics-key."+o.AppName+".log")
	lf, err := rotate.NewFile(f, o.MaxBackups)

	if err != nil {
		logrus.Warnf("fail to new log file %s", f)
	}

	return &Runner{
		Interval: o.Interval,
		C:        make(chan Line, o.ChanCap),
		stop:     make(chan bool, 1),
		Logfile:  lf,
	}
}

// Start starts the runner
func (r *Runner) Start() {
	go r.run()
	runtime.SetFinalizer(r, func(r *Runner) { r.Stop() })
}

// Stop stops the runner
func (r *Runner) Stop() {
	select {
	case r.stop <- true:
	default:
	}
}

func (r *Runner) run() {
	r.startTime = time.Now()

	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()

	for {
		select {
		case ml := <-r.C:
			jsonLog := util.JSONCompact(ml)

			r.writeLog(jsonLog)

			if r.afterInterval() {
				r.log()
			}
		case <-ticker.C:
			fmt.Println(time.Now())
			r.log()
		case <-r.stop:
			logrus.Info("runner stopped")
			return
		}
	}
}

func (r *Runner) afterInterval() bool {
	return time.Since(r.startTime) > r.Interval
}

func (r *Runner) log() {
	r.startTime = time.Now()
}

func (r *Runner) writeLog(content string) {
	if r.Logfile == nil {
		return
	}

	if _, err := r.Logfile.Write([]byte(content + "\n")); err != nil {
		logrus.Warnf("fail to write log, error %v", err)
	}
}
