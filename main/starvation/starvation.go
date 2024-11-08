package starvation

import "fmt"
import "time"

func snd(ch chan int) {
	var x int = 0
	for {
		x++
		ch <- x
		time.Sleep(1 * 1e9)
	}

}

func rcv(ch chan int) {
	var x int
	for {
		x = <-ch
		fmt.Printf("received %d \n", x)
	}

}

func starvation() {
	var ch chan int = make(chan int)
	go rcv(ch) // R1
	go snd(ch) // S1
	rcv(ch)    // R2

}
