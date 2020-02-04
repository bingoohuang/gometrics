package gometrics_test

import (
	"testing"

	"github.com/bingoohuang/gometrics"
)

func TestHitRate(t *testing.T) {
	fr := gometrics.HitRate("key1", "key2", "key3")
	fr.IncrHit()
	fr.IncrTotal()
}
