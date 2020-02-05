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
		AppName:         "bingoohuangapp",
		MetricsInterval: 3 * time.Second,
		HBInterval:      20 * time.Second,
		ChanCap:         123,
		LogPath:         "/tmp/metricslog",
		MaxBackups:      7,
	}, o)
}
