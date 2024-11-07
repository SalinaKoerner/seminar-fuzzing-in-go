package main

import "testing"

func FuzzRunner(f *testing.F) {
	testinputs := []int{5}

	for _, tc := range testinputs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, in int) {
		Runner(in)
	})
}
