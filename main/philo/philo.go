package main

import "fmt"
import "time"

// N philosophers sit at a table with a total of N forks.
// Each philosopher requires 2 forks to eat.

func philo(id int, forks chan int) {

	for {
		// take 2 forks
		<-forks
		<-forks
		fmt.Printf("%d eats \n", id)
		time.Sleep(1 * 1e9)
		// put back 2 forks
		forks <- 1
		forks <- 1

		time.Sleep(1 * 1e9) // think
	}

}

func main() {
	var forks = make(chan int, 3)
	forks <- 1 // put 3 forks on the table
	forks <- 1
	forks <- 1
	go philo(1, forks)
	go philo(2, forks)
	philo(3, forks)
}
