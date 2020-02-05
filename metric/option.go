package metric

import (
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/creasty/defaults"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func EnvOption(filenames ...string) *Option {
	envFile := os.Getenv("ENV_FILE")
	if len(filenames) == 0 && envFile != "" {
		filenames = []string{envFile}
	}

	if err := godotenv.Load(filenames...); err != nil {
		logrus.Warnf("loading env file error %+v", err)
	}

	o := &Option{}
	if err := env.Parse(o); err != nil {
		logrus.Warnf("parse env to option error %+v", err)
	}

	setDefaults(o)

	return o
}

// Option defines the option for runner
type Option struct {
	AppName    string        `env:"APP_NAME"`                            // 应用名称，默认使用当前进程的PID
	Interval   time.Duration `default:"1s" env:"INTERVAL"`               // 每隔多少时间记录一次日志
	ChanCap    int           `default:"1000" env:"CHAN_CAP"`             // 指标通道容量，当指标大量发送容量堆满时，自动扔弃
	LogPath    string        `default:"/var/log/metrics" env:"LOG_PATH"` // 日志路径
	MaxBackups int           `default:"7" env:"MAX_BACKUPS"`             // 最大保留天数
}

// OptionFn defines the function for options setting
type OptionFn func(o *Option)

// AppName sets the app name
func AppName(v string) OptionFn { return func(o *Option) { o.AppName = v } }

// Interval sets the interval to write log line
func Interval(v time.Duration) OptionFn { return func(o *Option) { o.Interval = v } }

// ChanCap sets the capacity of metrics line channel
func ChanCap(v int) OptionFn { return func(o *Option) { o.ChanCap = v } }

// LogPath sets the log path of metrics log files
func LogPath(v string) OptionFn { return func(o *Option) { o.LogPath = v } }

// MaxBackups sets max backups of metrics log files
func MaxBackups(v int) OptionFn { return func(o *Option) { o.MaxBackups = v } }

// NewRunner creates a Runner
func NewRunner(ofs ...OptionFn) *Runner {
	return NewRunnerOptions(createOption(ofs))
}

func createOption(ofs []OptionFn) *Option {
	o := &Option{}

	for _, of := range ofs {
		of(o)
	}

	setDefaults(o)

	return o
}

func setDefaults(o *Option) {
	if err := defaults.Set(o); err != nil {
		logrus.Warnf("defaults set error %v", err)
	}

	if o.AppName == "" {
		o.AppName = strconv.Itoa(os.Getpid())
	}
}
