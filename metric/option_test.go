package metric

import (
	"os"
	"testing"
	"time"

	"github.com/bingoohuang/gometrics/util"
	"github.com/stretchr/testify/assert"
)

// nolint gomnd
func TestEnvOption(t *testing.T) {
	println(util.PickFirst(os.Getwd()))
	var o Option
	EnvOption("../testdata/golden.env")(&o)
	assert.Equal(t, Option{
		LogFilePrefix:  "metrics-key.",
		LogFilePostfix: ".log",
		AppName:        "bingoohuangapp",
		Interval:       3 * time.Second,
		ChanCap:        123,
		LogPath:        "/tmp/metricslog",
		MaxBackups:     7,
	}, o)
}
