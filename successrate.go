package gometrics

// SuccessRateRecorder record success rate
type SuccessRateRecorder struct {
	Recorder
}

// SuccessRate makes a SuccessRateRecorder
func SuccessRate(keys ...string) SuccessRateRecorder {
	return DefaultRunner.SuccessRate(keys...)
}

// SuccessRate makes a SuccessRateRecorder
func (r *Runner) SuccessRate(keys ...string) SuccessRateRecorder {
	return SuccessRateRecorder{r.MakeRecorder("SUCCESS_RATE", keys)}
}

// IncrSuccess increment success count
func (c SuccessRateRecorder) IncrSuccess() {
	c.PutRecord(1, 0, 0, 0) // nolint gomnd
}

// IncrTotal increment total
func (c SuccessRateRecorder) IncrTotal() {
	c.PutRecord(0, 1, 0, 0) // nolint gomnd
}
