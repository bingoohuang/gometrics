package gometrics_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/bingoohuang/gometrics"
)

// nolint gomnd
func TestRTRecorder_Record(t *testing.T) {
	rt1 := gometrics.RT("key1")
	rt2 := gometrics.RT("key1", "key2")
	rt3 := gometrics.RT("key1", "key2", "key3")

	f := func() {
		time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
	}

	c := make(chan bool)

	go func() {
		f()
		rt1.Record()
		c <- true
	}()

	go func() {
		f()
		rt2.Record()
		c <- true
	}()

	go func() {
		f()
		rt3.Record()
		c <- true
	}()

	<-c
	<-c
	<-c
}

// nolint gomnd
func TestQPS_Record(t *testing.T) {
	gometrics.QPS("key1").Record(1)
	gometrics.QPS("key1", "key2").Record(1)
	gometrics.QPS("key1", "key2", "key3").Record(1)
}

func TestSuccessRate(t *testing.T) {
	sr := gometrics.SuccessRate("key1", "key2", "key3")
	sr.IncrSuccess()
	sr.IncrTotal()
}

func TestFailRate(t *testing.T) {
	fr := gometrics.FailRate("key1", "key2", "key3")
	fr.IncrFail()
	fr.IncrTotal()
}

func TestHitRate(t *testing.T) {
	fr := gometrics.HitRate("key1", "key2", "key3")
	fr.IncrHit()
	fr.IncrTotal()
}

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
