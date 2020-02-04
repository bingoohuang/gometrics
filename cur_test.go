package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
)

// nolint gomnd
func TestCurRecorder_Record(t *testing.T) {
	c1 := gometrics.Cur("key1")
	c2 := gometrics.Cur("key1", "key2")
	c3 := gometrics.Cur("key1", "key2", "key3")

	c1.Record(1)
	c2.Record(2)
	c3.Record(3)

	c1.Record(4)
	c2.Record(5)
	c3.Record(6)
}
