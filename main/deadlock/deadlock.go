package deadlock

import "fmt"

func snd(ch chan int) {
	var x int = 0
	x++
	ch <- x
}

func rcv(ch chan int) {
	var x int
	x = <-ch
	fmt.Printf("received %d \n", x)

}

func deadlock() {
	var ch chan int = make(chan int)
	go rcv(ch) // R1
	go snd(ch) // S1
	rcv(ch)    // R2

}
