package main

import (
	"fmt"
	"testing"
)

func FuzzFunny(f *testing.F) {
	testinputs := []int{5, -334, 4958}
	testinputs2 := []int{454, 222, -1}

	for _, tc := range testinputs {
		for _, tc2 := range testinputs2 {
			f.Add(tc, tc2) // Use f.Add to provide a seed corpus
		}

	}
	f.Fuzz(func(t *testing.T, in int, in2 int) {
		fmt.Printf("%d => %d \n", in, in2)
		Funny(in, in2)

	})
}
