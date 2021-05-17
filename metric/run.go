package metric

import (
	"io"
	"math"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/bingoohuang/gometrics/pkg/rotate"
	"github.com/bingoohuang/gometrics/pkg/util"
	"github.com/sirupsen/logrus"
)

// DefaultRunner is the default runner for metric recording.
var DefaultRunner = NewRunner(EnvOption())

// Stop stops the default runner.
func Stop() {
	DefaultRunner.Stop()
}

// Runner is a runner for metric rotate writing.
type Runner struct {
	startTime time.Time
	AppName   string

	MetricsInterval time.Duration
	HBInterval      time.Duration

	C    chan Line
	stop chan bool

	MetricsLogfile io.Writer
	HBLogfile      io.Writer

	cache     map[cacheKey]*Line
	option    *Option
	startOnce sync.Once
}

type cacheKey struct {
	Key     string
	LogType LogType
}

func makeCacheKey(key string, logType LogType) cacheKey {
	return cacheKey{Key: key, LogType: logType}
}

// NewRunner creates a Runner.
func NewRunner(ofs ...OptionFn) *Runner {
	o := CreateOption(ofs...)

	r := &Runner{
		option:          o,
		AppName:         o.AppName,
		MetricsInterval: o.MetricsInterval,
		HBInterval:      o.HBInterval,
		C:               make(chan Line, o.ChanCap),
		stop:            make(chan bool, 1),
		cache:           make(map[cacheKey]*Line),
	}

	runtime.SetFinalizer(r, func(r *Runner) { r.Stop() })

	return r
}

func createRotateFile(o *Option, prefix string) io.Writer {
	f := filepath.Join(o.LogPath, prefix+o.AppName+".log")
	lf, err := rotate.NewFile(f, o.MaxBackups)
	if err != nil {
		logrus.Warnf("fail to new logMetrics file %s, error %v", f, err)
		return nil
	}

	return lf
}

// Start starts the runner.
func (r *Runner) start() {
	o := r.option
	r.MetricsLogfile = createRotateFile(o, "metrics-key.")
	r.HBLogfile = createRotateFile(o, "metrics-hb.")

	go r.run()

	logrus.Info("runner started")
}

// Stop stops the runner.
func (r *Runner) Stop() {
	select {
	case r.stop <- true:
	default:
	}
}

func (r *Runner) run() {
	r.startTime = time.Now()

	metricsTicker := time.NewTicker(r.MetricsInterval)
	defer metricsTicker.Stop()

	r.logHB()

	hbTicker := time.NewTicker(r.HBInterval)
	defer hbTicker.Stop()

	for {
		select {
		case l := <-r.C:
			r.mergeLog(l)

			if r.afterMetricsInterval() {
				r.logMetrics()
			}
		case <-metricsTicker.C:
			r.logMetrics()
		case <-hbTicker.C:
			r.logHB()
		case <-r.stop:
			logrus.Info("runner stopped")

			return
		}
	}
}

func (r *Runner) afterMetricsInterval() bool { return time.Since(r.startTime) > r.MetricsInterval }

func (r *Runner) logMetrics() {
	if r.MetricsLogfile == nil {
		return
	}

	r.startTime = time.Now()

	for k, pv := range r.cache {
		v := *pv

		// 处理瞬间current > total的情况.
		if v.LogType.isPercent() && v.V1 > v.V2 {
			v.V1 = v.V2
		}

		if FloatEquals(v.V1, 0) && FloatEquals(v.V2, 0) {
			continue
		}

		r.writeLog(r.MetricsLogfile, v)

		if v.LogType.isSimple() {
			delete(r.cache, k)
		} else {
			pv.V1 -= v.V1
			pv.V2 -= v.V2
		}
	}
}

func (r *Runner) writeLog(file io.Writer, v Line) {
	v.Time = time.Now().Format(TimeLayout)
	v.Hostname = util.Hostname

	if r.option.Debug {
		s, _ := v.ToLineProtocol()
		logrus.Infof("LineProtocol:%s", s)
	}
	if file == nil {
		return
	}

	content := util.JSONCompact(v)

	if _, err := file.Write([]byte(content + "\n")); err != nil {
		logrus.Warnf("fail to write log of metrics, error %+v", err)
	}
}

func (r *Runner) mergeLog(l Line) {
	k := makeCacheKey(l.Key, l.LogType)
	if c, ok := r.cache[k]; ok {
		if l.LogType.isSimple() { // 瞬值，直接更新日志
			c.V1 = l.V1
			c.V2 = l.V2
		}

		c.updateMinMax(l)
	} else {
		r.cache[k] = &l
	}
}

func (r *Runner) logHB() {
	if r.HBLogfile == nil {
		return
	}

	r.writeLog(r.HBLogfile, Line{
		Key:     r.AppName + ".hb",
		LogType: HB,
		V1:      1,
	})
}

func (l *Line) updateMinMax(n Line) {
	uv1, uv2, curMin, curMax := l.V1+n.V1, l.V2+n.V2, l.Min, l.Max

	// 百分比类型时，uv1 > uv2没意义（可能是分子还没更新，分母累积提前到达）
	if n.V2 <= EPSILON || n.LogType.isPercent() && uv1 > uv2 {
		l.update(uv1, uv2, curMin, curMax)

		return
	}

	var ratio float64 = 1
	if n.LogType.isPercent() {
		ratio = 100
	}

	if n.LogType.isUseCurrent4MinMax() {
		ratio *= n.V1 / n.V2
	} else {
		ratio *= uv1 / uv2
	}

	l.update(uv1, uv2, Min(curMin, ratio), Max(curMax, ratio))
}

func (l *Line) update(v1, v2, min, max float64) {
	l.V1 = v1
	l.V2 = v2
	l.Min = min
	l.Max = max
}

// Max returns the max of two number.
func Max(max, v float64) float64 {
	if max < EPSILON || v > max {
		return v
	}

	return max
}

// Min returns the min of two number.
func Min(min, v float64) float64 {
	if min < EPSILON || v < min {
		return v
	}

	return min
}

var EPSILON = math.Nextafter(1, 2) - 1

func FloatEquals(a, b float64) bool {
	return a-b < EPSILON && b-a < EPSILON
}
