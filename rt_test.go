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
