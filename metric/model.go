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
	V1       float64 `json:"v1"` // 小数
	V2       float64 `json:"v2"` // 只有比率类型的时候，才用到v2
	V3       float64 `json:"v3"` // RT 当 [300-400) ms 时 v3 = 1
	V4       float64 `json:"v4"` // RT 当 [400-500) ms 时 v4 = 1
	V5       float64 `json:"v5"` // RT 当 [500-600) ms 时 v5 = 1
	V6       float64 `json:"v6"` // RT 当 [600-700) ms 时 v6 = 1
	V7       float64 `json:"v7"` // RT 当 [700-800) ms 时 v7 = 1
	V8       float64 `json:"v8"` // RT 当 [800-900) ms 时 v8 = 1
	V9       float64 `json:"v9"` // RT 当 [900-∞) ms 时 v9 = 1

	Min float64 `json:"min"` // 累计最小值
	Max float64 `json:"max"` // 累计最大值
}

// ToLineProtocol print l to a influxdb v1 line protocol format.
func (l Line) ToLineProtocol() (string, error) {
	t, err := time.Parse(TimeLayout, l.Time)
	if err != nil {
		return "", err
	}

	return lineprotocol.Build(string(l.LogType),
		map[string]string{"key": l.Key, "hostname": l.Hostname},
		map[string]interface{}{"v1": l.V1, "v2": l.V2, "min": l.Min, "max": l.Max,
			"v3": l.V3, "v4": l.V4, "v5": l.V5, "v6": l.V6, "v7": l.V7, "v8": l.V8, "v9": l.V9},
		t)
}

// AsyncPut new a metric line.
func (r *Runner) AsyncPut(keys []string, logType LogType, v1, v2 float64, vx ...float64) {
	fv := func(idx int) float64 {
		if len(vx) > idx-3 {
			return vx[idx-3]
		}

		return 0
	}

	select {
	case r.C <- Line{
		Key:     strings.Join(keys, "#"),
		LogType: logType,
		V1:      v1,
		V2:      v2,
		Min:     -1,
		Max:     -1,
		V3:      fv(3),
		V4:      fv(4),
		V5:      fv(5),
		V6:      fv(6),
		V7:      fv(7),
		V8:      fv(8),
		V9:      fv(9),
	}: // processed already.
	default: // bypass, async.
	}
}
