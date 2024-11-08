package livelock

import (
	"testing"
	"time"
)

func FuzzLivelock(f *testing.F) {
	f.Add(1)

	f.Fuzz(func(t *testing.T, input int) {
		go livelock()
		time.Sleep(1 * time.Second)
	})
}
