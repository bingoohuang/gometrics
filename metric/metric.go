package metric

import "time"

// Recorder record rate
type Recorder struct {
	Runner  *Runner
	LogType LogType
	Rate    bool
	Keys
}

// MakeRecorder creates a Recorder
func MakeRecorder(logType LogType, keys []string) Recorder {
	return DefaultRunner.MakeRecorder(logType, keys)
}

// MakeRecorder creates a Recorder
func (r *Runner) MakeRecorder(logType LogType, keys []string) Recorder {
	return Recorder{Runner: r, LogType: logType, Keys: NewKeys(keys)}
}

// PutRecord put a metric record to channel
func (c Recorder) PutRecord(v1, v2, min, max int64) {
	if c.Checked {
		c.Runner.PutMetricLine(c.Keys.Keys, c.LogType, v1, v2, min, max)
	}
}

// RTRecorder is a Round-Time recorder 平均响应时间
type RTRecorder struct {
	Recorder
	Start time.Time
}

// RT makes a RT Recorder
func RT(keys ...string) RTRecorder { return DefaultRunner.RT(keys...) }

// RT makes a RT Recorder
func (r *Runner) RT(keys ...string) RTRecorder {
	return RTRecorder{Recorder: r.MakeRecorder(KeyRT, keys), Start: time.Now()}
}

// Record records a round-time
func (r RTRecorder) Record() { r.RecordSince(r.Start) }

// RecordSince records a round-time since start
func (r RTRecorder) RecordSince(start time.Time) {
	r.PutRecord(time.Since(start).Milliseconds(), 1, 0, 0)
}

// QPSRecorder is a QPS recorder
type QPSRecorder struct{ Recorder }

// QPS makes a QPS Recorder
func QPS(keys ...string) QPSRecorder { return DefaultRunner.QPS(keys...) }

// QPS makes a QPS Recorder
func (r *Runner) QPS(keys ...string) QPSRecorder { return QPSRecorder{r.MakeRecorder(KeyQPS, keys)} }

// Record records a request
func (q QPSRecorder) Record(times int64) {
	if times > 0 {
		q.PutRecord(times, 0, 0, 0)
	}
}

// SuccessRateRecorder record success rate
type SuccessRateRecorder struct{ Recorder }

// SuccessRate makes a SuccessRateRecorder
func SuccessRate(keys ...string) SuccessRateRecorder { return DefaultRunner.SuccessRate(keys...) }

// SuccessRate makes a SuccessRateRecorder
func (r *Runner) SuccessRate(keys ...string) SuccessRateRecorder {
	return SuccessRateRecorder{r.MakeRecorder(KeySuccessRate, keys)}
}

// IncrSuccess increment success count
func (c SuccessRateRecorder) IncrSuccess() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c SuccessRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}

// FailRateRecorder record success rate
type FailRateRecorder struct{ Recorder }

// FailRate creates a FailRateRecorder
func FailRate(keys ...string) FailRateRecorder { return DefaultRunner.FailRate(keys...) }

// FailRate creates a FailRateRecorder
func (r *Runner) FailRate(keys ...string) FailRateRecorder {
	return FailRateRecorder{r.MakeRecorder(KeyFailRate, keys)}
}

// IncrFail increment success count
func (c FailRateRecorder) IncrFail() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c FailRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}

// HitRateRecorder record hit rate
type HitRateRecorder struct{ Recorder }

// HitRate makes a HitRateRecorder
func HitRate(keys ...string) HitRateRecorder { return DefaultRunner.HitRate(keys...) }

// HitRate makes a HitRateRecorder
func (r *Runner) HitRate(keys ...string) HitRateRecorder {
	return HitRateRecorder{r.MakeRecorder(KeyHitRate, keys)}
}

// IncrHit increment success count
func (c HitRateRecorder) IncrHit() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c HitRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}

// CurRecorder record 瞬时值(Gauge)
type CurRecorder struct{ Recorder }

// Cur makes a Cur Recorder
func Cur(keys ...string) CurRecorder { return DefaultRunner.Cur(keys...) }

// Cur makes a Cur Recorder
func (r *Runner) Cur(keys ...string) CurRecorder {
	return CurRecorder{r.MakeRecorder(KeyCUR, keys)}
}

// Record record  v1
func (c CurRecorder) Record(v1 int64) {
	c.PutRecord(v1, 0, 0, 0)
}
