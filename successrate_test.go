package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
)

func TestSuccessRate(t *testing.T) {
	sr := gometrics.SuccessRate("key1", "key2", "key3")
	sr.IncrSuccess()
	sr.IncrTotal()
}
