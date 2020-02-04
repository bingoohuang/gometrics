package gometrics

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// DefaultRunner is the default runner for metric recording
var DefaultRunner = NewRunner(1*time.Second, 1000) // nolint

func init() { // nolint
	DefaultRunner.Start()
}

// Runner is a runner for metric log writing
type Runner struct {
	startTime time.Time
	Interval  time.Duration
	C         chan MetricLine
	StopC     chan bool
}

// NewRunner creates a Runner
func NewRunner(interval time.Duration, chanCap int) *Runner {
	return &Runner{
		Interval: interval,
		C:        make(chan MetricLine, chanCap),
		StopC:    make(chan bool, 1),
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
	case r.StopC <- true:
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
			fmt.Println(JSONCompact(ml))

			if r.afterInterval() {
				r.log()
			}
		case <-ticker.C:
			fmt.Println(time.Now())
			r.log()
		case <-r.StopC:
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
