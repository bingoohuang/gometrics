package metric

import (
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
)

// Option defines the option for runner
type Option struct {
	AppName         string        `env:"APP_NAME"`                            // 应用名称，默认使用当前进程的PID
	MetricsInterval time.Duration `default:"1s" env:"METRICS_INTERVAL"`       // 每隔多少时间记录一次日志
	HBInterval      time.Duration `default:"20s" env:"HB_INTERVAL"`           // 每隔多少时间记录一次心跳日志
	ChanCap         int           `default:"1000" env:"CHAN_CAP"`             // 指标通道容量，当指标大量发送容量堆满时，自动扔弃
	LogPath         string        `default:"/var/log/metrics" env:"LOG_PATH"` // 日志路径
	MaxBackups      int           `default:"7" env:"MAX_BACKUPS"`             // 最大保留天数
}

// OptionFn defines the function for options setting
type OptionFn func(o *Option)

// EnvOption sets the options from env
func EnvOption(filenames ...string) OptionFn {
	return func(o *Option) {
		envFile := os.Getenv("ENV_FILE")
		if len(filenames) == 0 && envFile != "" {
			filenames = []string{envFile}
		}

		if err := godotenv.Load(filenames...); err != nil {
			logrus.Warnf("loading env file error %+v", err)
		}

		if err := env.Parse(o); err != nil {
			logrus.Warnf("parse env to option error %+v", err)
		}

		o.SetDefaults()
	}
}

// AppName sets the app name
func AppName(v string) OptionFn { return func(o *Option) { o.AppName = v } }

// MetricsInterval sets the interval to write logMetrics line
func MetricsInterval(v time.Duration) OptionFn { return func(o *Option) { o.MetricsInterval = v } }

// ChanCap sets the capacity of metrics line channel
func ChanCap(v int) OptionFn { return func(o *Option) { o.ChanCap = v } }

// LogPath sets the logMetrics path of metrics logMetrics files
func LogPath(v string) OptionFn { return func(o *Option) { o.LogPath = v } }

// MaxBackups sets max backups of metrics logMetrics files
func MaxBackups(v int) OptionFn { return func(o *Option) { o.MaxBackups = v } }

// CreateOption creates Option by option functions
func CreateOption(ofs ...OptionFn) *Option {
	o := &Option{}

	for _, of := range ofs {
		of(o)
	}

	o.SetDefaults()

	return o
}

// SetDefaults set the default values to Option
func (o *Option) SetDefaults() {
	if err := defaults.Set(o); err != nil {
		logrus.Warnf("defaults set error %v", err)
	}

	if o.AppName == "" {
		o.AppName = strconv.Itoa(os.Getpid())
	}
}
