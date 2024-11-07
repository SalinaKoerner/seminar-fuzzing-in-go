package main

import "testing"

func FuzzPhilo(f *testing.F) {

	f.Add(1)
	f.Add(2)
	f.Add(3)

	forks := make(chan int, 3)
	forks <- 1
	forks <- 1
	forks <- 1

	f.Fuzz(func(t *testing.T, id int) {
		go philo(id, forks)
	})
}
