package main

// Deadlock possible!
func Runner(i int) {

	ch := make(chan int)

	go func() {
		ch <- i
	}()

	go func() {
		<-ch
	}()

	ch <- i

}

func main4() {

	for i := 0; i < 50; i++ {
		Runner(i)
	}

}
