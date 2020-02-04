package gometrics

// CurRecorder record 瞬时值(Gauge)
type CurRecorder struct {
	Recorder
}

// Cur makes a Cur Recorder
func Cur(keys ...string) CurRecorder {
	return DefaultRunner.Cur(keys...)
}

// Cur makes a Cur Recorder
func (r *Runner) Cur(keys ...string) CurRecorder {
	return CurRecorder{r.MakeRecorder("CUR", keys)}
}

// Record record  v1
func (c CurRecorder) Record(v1 int64) {
	c.PutRecord(v1, 0, 0, 0)
}
