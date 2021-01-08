package metric

import (
	"strings"
	"time"

	"github.com/bingoohuang/gometrics/pkg/lineprotocol"
)

// LogType means the logMetrics type.
type LogType string

const (
	// KeyRT RT 日志类型.
	KeyRT LogType = "RT"
	// KeyQPS QPS 日志类型.
	KeyQPS = "QPS"
	// KeySuccessRate SuccessRate 日志类型.
	KeySuccessRate = "SUCCESS_RATE"
	// KeyFailRate FailRate 日志类型.
	KeyFailRate = "FAIL_RATE"
	// KeyHitRate HitRate 日志类型.
	KeyHitRate = "HIT_RATE"
	// KeyCUR CUR 日志类型.
	KeyCUR = "CUR"

	// HB 特殊处理，每?s记录一次.
	HB = "HB"
)

// isSimple 是否简单的值，值与值之间，不需要有累计等关系.
func (lt LogType) isSimple() bool { return lt == KeyCUR }

// isUseCurrent4MinMax 是否使用当前v1/v2值来生成，还是使用累积值来生成min/max值.
func (lt LogType) isUseCurrent4MinMax() bool { return lt == KeyRT }

// isPercent 是否是百分比类型.
func (lt LogType) isPercent() bool {
	switch lt {
	case KeySuccessRate, KeyFailRate, KeyHitRate:
		return true
	default:
		return false
	}
}

const TimeLayout = "20060102150405000"

// Line represents a metric rotate line structure in rotate file.
type Line struct {
	Time     string  `json:"time"` // yyyyMMddHHmmssSSS
	Key      string  `json:"key"`  // {{k1}}#{{k2}}#{{k3}}
	Hostname string  `json:"hostname"`
	LogType  LogType `json:"logtype"`
	V1       float64 `json:"v1"`  // 小数
	V2       float64 `json:"v2"`  // 只有比率类型的时候，才用到v2
	Min      float64 `json:"min"` // 累计最小值
	Max      float64 `json:"max"` // 累计最大值
}

// ToLineProtocol
func (l Line) ToLineProtocol() (string, error) {
	t, err := time.Parse(TimeLayout, l.Time)
	if err != nil {
		return "", err
	}

	return lineprotocol.Build(string(l.LogType),
		map[string]string{"key": l.Key, "hostname": l.Hostname},
		map[string]interface{}{"v1": l.V1, "v2": l.V2, "min": l.Min, "max": l.Max},
		t)
}

// AsyncPut new a metric line.
func (r *Runner) AsyncPut(keys []string, logType LogType, v1, v2 float64) {
	r.startOnce.Do(r.start)

	select {
	case r.C <- Line{
		Key:     strings.Join(keys, "#"),
		LogType: logType,
		V1:      v1,
		V2:      v2,
		Min:     -1,
		Max:     -1,
	}: // processed already.
	default: // bypass, async.
	}
}
