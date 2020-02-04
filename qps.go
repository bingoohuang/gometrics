package gometrics

// QPSRecorder is a QPS recorder
type QPSRecorder struct {
	Recorder
}

// QPS makes a QPS Recorder
func QPS(keys ...string) QPSRecorder {
	return DefaultRunner.QPS(keys...)
}

// QPS makes a QPS Recorder
func (r *Runner) QPS(keys ...string) QPSRecorder {
	return QPSRecorder{r.MakeRecorder("QPS", keys)}
}

// Record records a request
func (q QPSRecorder) Record(times int64) {
	if times > 0 {
		q.PutRecord(times, 0, 0, 0)
	}
}
