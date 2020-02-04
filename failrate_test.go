package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
)

func TestFailRate(t *testing.T) {
	fr := gometrics.FailRate("key1", "key2", "key3")
	fr.IncrFail()
	fr.IncrTotal()
}
