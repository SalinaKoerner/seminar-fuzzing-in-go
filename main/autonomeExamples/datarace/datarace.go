package datarace

import "fmt"
import "time"

func datarace(delay time.Duration) {
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

	time.Sleep(delay)
	fmt.Printf("done \n")

}
