package gometrics

// FailRateRecorder record success rate
type FailRateRecorder struct {
	Recorder
}

// FailRate creates a FailRateRecorder
func FailRate(keys ...string) FailRateRecorder {
	return DefaultRunner.FailRate(keys...)
}

// FailRate creates a FailRateRecorder
func (r *Runner) FailRate(keys ...string) FailRateRecorder {
	return FailRateRecorder{r.MakeRecorder("FAIL_RATE", keys)}
}

// IncrFail increment success count
func (c FailRateRecorder) IncrFail() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c FailRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}
