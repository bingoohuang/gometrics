package main

import (
	"math/rand"
	"time"

	"github.com/bingoohuang/golog"

	"github.com/bingoohuang/gometrics/metric"
)

func init() {
	_ = golog.Setup()
	metric.DefaultRunner.Stop()
	metric.DefaultRunner = metric.NewRunner(metric.EnvOption())
	metric.DefaultRunner.Start()
}

func main() {
	f := func() {
		time.Sleep(100 + time.Duration(rand.Int31n(900))*time.Millisecond)
	}

	for i := 0; ; i++ {
		f()

		m := i % 6
		switch m {
		case 0:
			func() {
				defer metric.RT("key1", "key2", "key3").Record()
				f()
			}()
		case 1:
			func() {
				metric.QPS("key1", "key2", "key3").Record(1)
			}()
		case 2:
			func() {
				sr := metric.SuccessRate("key1", "key2", "key3")
				defer sr.IncrTotal()

				if rand.Intn(3) == 0 {
					sr.IncrSuccess()
				}
			}()
		case 3:
			func() {
				fr := metric.FailRate("key1", "key2", "key3")
				defer fr.IncrTotal()

				if rand.Intn(10) == 0 {
					fr.IncrFail()
				}
			}()
		case 4:
			func() {
				fr := metric.HitRate("key1", "key2", "key3")
				defer fr.IncrTotal()

				if rand.Intn(5) == 0 {
					fr.IncrHit()
				}
			}()
		case 5:
			func() {
				metric.Cur("key1", "key2", "key3").Record(100)
			}()
		}
	}
}
