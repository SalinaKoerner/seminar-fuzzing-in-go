package main

import "fmt"

// function takes a number and checks if it's even
func Even(i int) bool {
	// bug: numbers greater than 100 are never identified as even
	if i > 100 {
		return false
	}

	// returns true, if number is even
	if i%2 == 0 {
		return true
	}

	// return false in every other case
	return false
}

func main() {
	fmt.Printf("Are these numbers even? \n")
	fmt.Printf("%d => %t \n", 5, Even(5))

	fmt.Printf("%d => %t \n", 0, Even(0))

}
