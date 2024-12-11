package main

import "testing"

type testPair struct {
	input    int
	expected bool
}

func TestEven(t *testing.T) {

	testcases := []testPair{
		{5, false}, {0, true}, {50, true}}

	for _, tc := range testcases {
		res := Even(tc.input)
		if res != tc.expected {
			t.Errorf("isEven: %d => %t, want %t", tc.input, res, tc.expected)
		}

	}

}

// Fuzz test for the Even function
func FuzzEven(f *testing.F) {

	// initial set of seed inputs
	testinputs := []int{5, 0, 50}

	for _, ti := range testinputs {
		f.Add(ti) // Use f.Add to provide a seed corpus
	}

	// fuzz target, with fuzzing argument "input"
	f.Fuzz(func(t *testing.T, input int) {
		result := Even(input)
		result2 := Even(input + 1)
		if result == result2 {
			// throws an error, if number and its successor both identified as even
			t.Errorf("Fail: %d => %t, %d => %t", input, result, input+1, result2)
		}
	})
}
