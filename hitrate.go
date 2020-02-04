package gometrics

// HitRateRecorder record hit rate
type HitRateRecorder struct {
	Recorder
}

// HitRate makes a HitRateRecorder
func HitRate(keys ...string) HitRateRecorder {
	return DefaultRunner.HitRate(keys...)
}

// HitRate makes a HitRateRecorder
func (r *Runner) HitRate(keys ...string) HitRateRecorder {
	return HitRateRecorder{r.MakeRecorder("HIT_RATE", keys)}
}

// IncrHit increment success count
func (c HitRateRecorder) IncrHit() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c HitRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}
