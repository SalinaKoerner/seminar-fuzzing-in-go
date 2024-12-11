package livelock

import (
	"fmt"
	"time"
)

func livelock() {
	var x int
	y := make(chan int, 1)

	// T2
	go func() {
		y <- 1
		x++
		<-y

	}()

	x++
	y <- 1
	<-y

	time.Sleep(1 * 1e9)
	fmt.Printf("done \n")

}
