package metric

import (
	"os"
	"strings"
)

// LogType means the logMetrics type
type LogType string

const (
	KeyRT          LogType = "RT"
	KeyQPS         LogType = "QPS"
	KeySuccessRate LogType = "SUCCESS_RATE"
	KeyFailRate    LogType = "FAIL_RATE"
	KeyHitRate     LogType = "HIT_RATE"
	KeyCUR         LogType = "CUR"

	// HB 特殊处理，每20s记录一次
	HB LogType = "HB"
)

// isSimple 是否简单的值，值与值之间，不需要有累计等关系
func (lt LogType) isSimple() bool { return lt == KeyCUR }

// isPercent 是否是百分比类型
func (lt LogType) isPercent() bool {
	switch lt {
	case KeySuccessRate, KeyFailRate, KeyHitRate:
		return true
	}

	return false
}

// isUseCurrentValue4MinMax 是否使用当前v1/v2值来生成，还是使用累积值来生成min/max值
func (lt LogType) isUseCurrentValue4MinMax() bool { return lt == KeyRT }

// Line represents a metric rotate line structure in rotate file
type Line struct {
	Time     string  `json:"time"` // yyyyMMddHHmmssSSS
	Key      string  `json:"key"`  // {{k1}}#{{k2}}#{{k3}}
	Hostname string  `json:"hostname"`
	LogType  LogType `json:"logtype"`
	V1       int64   `json:"v1"`
	V2       int64   `json:"v2"`
	Min      int64   `json:"min"`
	Max      int64   `json:"max"`
}

// Hostname stores hostname
var Hostname string // nolint

func init() { // nolint
	Hostname, _ = os.Hostname()
}

// PutMetricLine new a Line
func (r *Runner) PutMetricLine(keys []string, logType LogType, v1, v2, min, max int64) {
	select {
	case r.C <- Line{
		Key:      strings.Join(keys, "#"),
		Hostname: Hostname,
		LogType:  logType,
		V1:       v1,
		V2:       v2,
		Min:      min,
		Max:      max,
	}: // processed already
	default: // bypass
	}
}
