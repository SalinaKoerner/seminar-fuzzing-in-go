package datarace

import "testing"
import "time"

func FuzzDatarace(f *testing.F) {
	f.Add(50)

	f.Fuzz(func(t *testing.T, delayMs int) {
		delay := time.Duration(delayMs) * time.Millisecond
		datarace(delay)
	})
}
