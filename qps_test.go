package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
)

// nolint gomnd
func TestQPS_Record(t *testing.T) {
	gometrics.QPS("key1").Record(1)
	gometrics.QPS("key1", "key2").Record(1)
	gometrics.QPS("key1", "key2", "key3").Record(1)
}
