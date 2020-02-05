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
	o := EnvOption("../testdata/golden.env")
	assert.Equal(t, Option{
		AppName:    "bingoohuangapp",
		Interval:   3 * time.Second,
		ChanCap:    123,
		LogPath:    "/tmp/metricslog",
		MaxBackups: 7,
	}, *o)
}
