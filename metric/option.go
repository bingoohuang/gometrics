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
	LogFilePrefix  string        `default:"metrics-key."`                    // 日志文件前缀，例如metrics-key.
	AppName        string        `env:"APP_NAME"`                            // 应用名称，默认使用当前进程的PID
	LogFilePostfix string        `default:".log"`                            // 日志文件后缀，eg.log, .hb
	Interval       time.Duration `default:"1s" env:"INTERVAL"`               // 每隔多少时间记录一次日志
	ChanCap        int           `default:"1000" env:"CHAN_CAP"`             // 指标通道容量，当指标大量发送容量堆满时，自动扔弃
	LogPath        string        `default:"/var/log/metrics" env:"LOG_PATH"` // 日志路径
	MaxBackups     int           `default:"7" env:"MAX_BACKUPS"`             // 最大保留天数
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

		setDefaults(o)
	}
}

// AppName sets the app name
func AppName(v string) OptionFn { return func(o *Option) { o.AppName = v } }

// LogFilePrefix sets the log file name's prefix
func LogFilePrefix(v string) OptionFn { return func(o *Option) { o.LogFilePrefix = v } }

// LogFilePrefix sets the log file name's postfix
func LogFilePostfix(v string) OptionFn { return func(o *Option) { o.LogFilePostfix = v } }

// Interval sets the interval to write log line
func Interval(v time.Duration) OptionFn { return func(o *Option) { o.Interval = v } }

// ChanCap sets the capacity of metrics line channel
func ChanCap(v int) OptionFn { return func(o *Option) { o.ChanCap = v } }

// LogPath sets the log path of metrics log files
func LogPath(v string) OptionFn { return func(o *Option) { o.LogPath = v } }

// MaxBackups sets max backups of metrics log files
func MaxBackups(v int) OptionFn { return func(o *Option) { o.MaxBackups = v } }

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
