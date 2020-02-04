package gometrics

// Recorder record rate
type Recorder struct {
	Runner  *Runner
	LogType string
	Rate    bool
	Keys
}

// MakeRecorder creates a Recorder
func MakeRecorder(logType string, keys []string) Recorder {
	return DefaultRunner.MakeRecorder(logType, keys)
}

// MakeRecorder creates a Recorder
func (r *Runner) MakeRecorder(logType string, keys []string) Recorder {
	return Recorder{Runner: r, LogType: logType, Keys: NewKeys(keys)}
}

// PutRecord put a metric record to channel
func (c Recorder) PutRecord(v1, v2, min, max int64) {
	if c.Checked {
		c.Runner.PutMetricLine(c.Keys.Keys, c.LogType, v1, v2, min, max)
	}
}
