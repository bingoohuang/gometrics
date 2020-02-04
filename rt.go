package gometrics

import "time"

// RTRecorder is a Round-Time recorder 平均响应时间
type RTRecorder struct {
	Recorder
	Start time.Time
}

// RT makes a RT Recorder
func RT(keys ...string) RTRecorder {
	return DefaultRunner.RT(keys...)
}

// RT makes a RT Recorder
func (r *Runner) RT(keys ...string) RTRecorder {
	return RTRecorder{Recorder: r.MakeRecorder("RT", keys), Start: time.Now()}
}

// Record records a round-time
func (r RTRecorder) Record() {
	r.RecordSince(r.Start)
}

// RecordSince records a round-time since start
func (r RTRecorder) RecordSince(start time.Time) {
	r.PutRecord(time.Since(start).Milliseconds(), 1, 0, 0)
}
