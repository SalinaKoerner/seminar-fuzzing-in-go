package starvation

import (
	"testing"
)

func FuzzStarvation(f *testing.F) {
	f.Add(0)

	f.Fuzz(func(t *testing.T, _ int) {
		var ch chan int = make(chan int)

		go rcv(ch) // R1
		go snd(ch) // S1
		rcv(ch)    // R2

	})
}
