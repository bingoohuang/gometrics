package gometrics

import (
	"os"
	"strings"
	"time"
)

// MetricLine represents a metric rotate line structure in rotate file
type MetricLine struct {
	Time     string `json:"time"` // yyyyMMddHHmmssSSS
	Key      string `json:"key"`  // {{k1}}#{{k2}}#{{k3}}
	Hostname string `json:"hostname"`
	LogType  string `json:"logtype"`
	V1       int64  `json:"v1"`
	V2       int64  `json:"v2"`
	Min      int64  `json:"min"`
	Max      int64  `json:"max"`
}

// Hostname stores hostname
var Hostname string // nolint

func init() { // nolint
	Hostname, _ = os.Hostname()
}

// PutMetricLine new a MetricLine
func (r *Runner) PutMetricLine(keys []string, logtype string, v1, v2, min, max int64) {
	select {
	case r.C <- MetricLine{
		Time:     time.Now().Format("20060102150405000"),
		Key:      strings.Join(keys, "#"),
		Hostname: Hostname,
		LogType:  logtype,
		V1:       v1,
		V2:       v2,
		Min:      min,
		Max:      max,
	}: // processed already
	default: // bypass
	}
}
