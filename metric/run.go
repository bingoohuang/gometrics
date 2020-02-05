package metric

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/bingoohuang/gometrics/util"

	"github.com/bingoohuang/gometrics/rotate"
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

	cache map[cacheKey]*Line
}

type cacheKey struct {
	Key     string
	LogType LogType
}

func makeCacheKey(key string, logType LogType) cacheKey {
	return cacheKey{Key: key, LogType: logType}
}

// NewRunner creates a Runner
func NewRunner(ofs ...OptionFn) *Runner {
	return NewRunnerOptions(createOption(ofs))
}

// NewRunnerOptions creates a Runner
func NewRunnerOptions(o *Option) *Runner {
	f := filepath.Join(o.LogPath, "metrics-Key."+o.AppName+".log")
	lf, err := rotate.NewFile(f, o.MaxBackups)

	if err != nil {
		logrus.Warnf("fail to new log file %s", f)
	}

	return &Runner{
		Interval: o.Interval,
		C:        make(chan Line, o.ChanCap),
		stop:     make(chan bool, 1),
		Logfile:  lf,
		cache:    make(map[cacheKey]*Line),
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
		case l := <-r.C:
			r.mergeLog(l)

			if r.afterInterval() {
				r.log()
			}
		case <-ticker.C:
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

	for k, pv := range r.cache {
		v := *pv

		if v.V1 == 0 && v.V2 == 0 {
			continue
		}

		// 处理瞬间current > total的情况
		if v.LogType.isPercent() && v.V1 > v.V2 {
			v.V1 = v.V2
		}

		v.Time = time.Now().Format("20060102150405000")
		r.writeLog(util.JSONCompact(v))

		if v.LogType.isSimple() {
			delete(r.cache, k)
		} else {
			pv.V1 -= v.V1
			pv.V2 -= v.V2
		}
	}
}

func (r *Runner) writeLog(content string) {
	if r.Logfile == nil {
		return
	}

	if _, err := r.Logfile.Write([]byte(content + "\n")); err != nil {
		logrus.Warnf("fail to write log, error %v", err)
	}
}

func (r *Runner) mergeLog(l Line) {
	k := makeCacheKey(l.Key, l.LogType)
	if c, ok := r.cache[k]; ok {
		if l.LogType == LogTypeCUR { // 瞬值，直接记录日志
			c.V1 = l.V1
			c.V2 = l.V2
		}

		c.updateMinMax(l)
	} else {
		r.cache[k] = &l
	}
}

func (l *Line) updateMinMax(n Line) {
	uv1 := l.V1 + n.V1
	uv2 := l.V2 + n.V2
	curMin := l.Min
	curMax := l.Max

	percentType := n.LogType.isPercent()

	// 百分比类型时，uv1 > uv2没意义（可能是分子还没更新，分母累积提前到达）
	if n.V2 <= 0 || percentType && uv1 > uv2 {
		l.V1 = uv1
		l.V2 = uv2
		l.Min = curMin
		l.Max = curMax

		return
	}

	var ratio int64 = 1
	if percentType {
		ratio = 100
	}

	if n.LogType.isUseCurrentValue4MinMax() {
		ratio *= n.V1 / n.V2
	} else {
		ratio *= uv1 / uv2
	}

	min := curMin
	if curMin < 0 || ratio < curMin {
		min = ratio
	}

	max := curMax
	if curMax < 0 || curMax < ratio {
		max = ratio
	}

	l.V1 = uv1
	l.V2 = uv2
	l.Min = min
	l.Max = max
}
