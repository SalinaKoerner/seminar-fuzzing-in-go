package deadlock

import "testing"

func FuzzDeadlock(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, input int) {
		var ch chan int = make(chan int)

		go rcv(ch) // R1
		go snd(ch) // S1
		rcv(ch)    //R2
	})
}
