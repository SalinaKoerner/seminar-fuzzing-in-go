package race

import "testing"

func FuzzRace(f *testing.F) {
	testinputs := []int{5}

	for _, tc := range testinputs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, in int) {
		Race(in)
	})
}
