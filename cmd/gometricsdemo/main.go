package main

import (
	"math/rand"
	"time"

	"github.com/bingoohuang/gometrics"
)

// nolint gomnd
func main() {
	f := func() {
		time.Sleep(100 + time.Duration(rand.Int31n(900))*time.Millisecond)
	}

	for i := 0; ; i++ {
		f()

		switch i % 6 {
		case 0:
			func() {
				defer gometrics.RT("key1", "key2", "key3").Record()
				f()
			}()
		case 1:
			func() {
				gometrics.QPS("key1", "key2", "key3").Record(1)
			}()
		case 2:
			func() {
				sr := gometrics.SuccessRate("key1", "key2", "key3")
				defer sr.IncrTotal()

				if rand.Intn(3) == 0 {
					sr.IncrSuccess()
				}
			}()
		case 3:
			func() {
				fr := gometrics.FailRate("key1", "key2", "key3")
				defer fr.IncrTotal()

				if rand.Intn(10) == 0 {
					fr.IncrFail()
				}
			}()
		case 4:
			func() {
				fr := gometrics.HitRate("key1", "key2", "key3")
				defer fr.IncrTotal()

				if rand.Intn(5) == 0 {
					fr.IncrHit()
				}
			}()
		case 5:
			func() {
				gometrics.Cur("key1", "key2", "key3").Record(100)
			}()
		}

	}
}
