F<!-- TOC -->
* [Fuzz Testing](#fuzz-testing)
* [Requirements](#requirements)
* [Fuzz test](#fuzz-test)
  * [The Output](#the-output)
  * [How the Go fuzzer works](#how-the-go-fuzzer-works)
* [General use cases of Go Fuzzer](#general-use-cases-of-go-fuzzer)
* [Data Race Detection with the Go Fuzzer](#data-race-detection-with-the-go-fuzzer)
  * [Data Race Example](#data-race-example)
* [Summary](#summary)
* [Sources / Literature](#sources--literature)
<!-- TOC -->

# Fuzz Testing

// todo: description of fuzz testing 


# Requirements

- Go 1.18 or later
- ARM64 or AMD64 architecture

# Fuzz test

Here you can see the structure of a fuzz test in Go.
The fuzz test needs to be in a file called {FileToTest}_test.go. The name of the fuzz test needs to start with "Fuzz"
followed by the name of the function, e.g. Fuzz{FunctionName}).

![img.png](img.png)

- Fuzz test: The fuzz test is the entire test with all its required parts, including fuzz target, fuzzing arguments and
  seed additions.
- Fuzz target: The function that is being executed with the corpus entries and the generated inputs.
- Fuzzing arguments: The fuzzing arguments are the data types that are being passed to the fuzzing function. These are
  also mutated to generate new inputs.
    - Seed corpus addition: That is the arguments provided by the programmer to be added to the corpus.

## The Output

```
λ go test -fuzz=FuzzDatarace -race
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 971 (324/sec), new interesting: 4 (total: 7)
fuzz: elapsed: 6s, execs: 1870 (300/sec), new interesting: 4 (total: 7)
fuzz: elapsed: 9s, execs: 2674 (268/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 12s, execs: 3643 (323/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 15s, execs: 4355 (237/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 16s, execs: 4400 (38/sec), new interesting: 5 (total: 8)
--- FAIL: FuzzDatarace (16.19s)
    --- FAIL: FuzzDatarace (0.00s)
        testing.go:1399: race detected during execution of test

    Failing input written to testdata\fuzz\FuzzDatarace\b5dd7ee3ea717225
    To re-run:
    go test -run=FuzzDatarace/b5dd7ee3ea717225
FAIL
exit status 1
FAIL    main/main/datarace      16.479s
```

- **-fuzz=FuzzDataRace:** the name of the fuzz test (in this case FuzzDataRace)
- **baseline coverage:** running the function with an initial set of inputs to ensure code coverage.
- **fuzzing with 8 workers:** means there are 8 concurrent test runners, this corresponds to the number of kernels of
  the pc
- **elapsed:** seconds since the fuzzing process started
- **execs:** number inputs tested, number of functions executed
- **new interesting:** number of newly discovered inputs that lead to previously undiscovered code paths.
  An input is considered interesting if it expands the code coverage with more than what the currently generated corpus
  already reached. In the brackets you can see the total size of the generated corpus.

- after over 2.618.761 executions the process failed, which means there is a problem in the function

## How the Go fuzzer works

When running a fuzz test Go starts with the seed corpus to get an initial generated corpus.
With a mutator it creates new inputs by generating new inputs, randomly modifying inputs or combining them.

The Go built-in fuzzer takes the provided inputs and runs it

The Go fuzzer if coverage guided. That means it learns from the code coverage that is expanded by new inputs. It tries
to explore as many new code paths as possible.

The "new interesting" amount will increase a lot in the beginning, but as the fuzzing process continues to work, there
will be less and less new interestings, as it won't find as many new code paths anymore.

It keeps track of the "new interesting" and tries to mutate those,so inputs hopefully lead to a previously undiscovered
code path.

The fuzzing arguments then are mutated to generate more inputs, that can be tested.

The fuzzing process stops when the fuzzer found a bug or the user manually stops it.
It can fail because of panics, runtime errors, asserting t.Error or when user written validation logic fails.

If the fuzzer finds a failing input it minimizes the input as small as possible. That input still has to reproduce the
bug. It's then added to the corpus and can be found in a directory usually called "testdata/fuzz/{FuzzTestName}".
Since it is save it can be used for future test runs.

# General use cases of Go Fuzzer

- Where deterministic unit tests don't discover bugs fuzz testing can help discover them by showing how the code behaves
  when encountering random inputs.
- Incorrect rune decoding
- Finding bugs when parsing inputs
- data races

# Data Race Detection with the Go Fuzzer

Here is an example that includes a data race. The two go routines (main and T2) both try to increment (a write
operation)
the variable x.

## Data Race Example

```
package main

import "fmt"
import "time"

func main() {
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
```

Fuzz test

```
package datarace

import "testing"
import "time"

func FuzzDatarace(f *testing.F) {
f.Add(50)

	f.Fuzz(func(t *testing.T, delayMs int) {
		delay := time.Duration(delayMs) * time.Millisecond
		datarace(delay)
	})
}
```

Sample runs

one of the first runs:

```
λ go test -fuzz=FuzzDatarace -race
fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 971 (324/sec), new interesting: 4 (total: 7)
fuzz: elapsed: 6s, execs: 1870 (300/sec), new interesting: 4 (total: 7)
fuzz: elapsed: 9s, execs: 2674 (268/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 12s, execs: 3643 (323/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 15s, execs: 4355 (237/sec), new interesting: 5 (total: 8)
fuzz: elapsed: 16s, execs: 4400 (38/sec), new interesting: 5 (total: 8)
--- FAIL: FuzzDatarace (16.19s)
    --- FAIL: FuzzDatarace (0.00s)
        testing.go:1399: race detected during execution of test

    Failing input written to testdata\fuzz\FuzzDatarace\b5dd7ee3ea717225
    To re-run:
    go test -run=FuzzDatarace/b5dd7ee3ea717225
FAIL
exit status 1
FAIL    main/main/datarace      16.479s
```

Running the fuzz test with the -race flag turns out to detect a race within around 17 seconds.

a later run:

```
λ go test -fuzz=FuzzDatarace -race
fuzz: elapsed: 0s, gathering baseline coverage: 0/19 completed
fuzz: elapsed: 0s, gathering baseline coverage: 19/19 completed, now fuzzing with 8 workers
fuzz: elapsed: 2s, execs: 982 (502/sec), new interesting: 0 (total: 19)
--- FAIL: FuzzDatarace (1.97s)
    --- FAIL: FuzzDatarace (0.00s)
        testing.go:1399: race detected during execution of test

    Failing input written to testdata\fuzz\FuzzDatarace\b45ecfe30afa5f31
    To re-run:
    go test -run=FuzzDatarace/b45ecfe30afa5f31
FAIL
exit status 1
FAIL    main/main/datarace      2.242s
```

The later runs reports a race much quicker than the first one, because it has more code coverage.

# Summary

To summarize one can say that the built-in Go fuzzer has proven to be a useful tool for finding some bugs.

However not all generated inputs are useful for testing, as some inputs would not occur in the day to day use of a
programm.
So it can be helpful to add a skip statement to eliminate the execution of unnecessary tests.

For data race the Go fuzzer seems to be better at detecting those, whereas Gos race detector often fails to find a race.

Although fuzz testing can discover unexpected bugs, it is still not a testing type that finds every possible bug in a
project. So fuzzing can prove the existence of a bug, but not the absence of one.

Go fuzzing in general also uses more resources than unit tests, because generating a lot of inputs and running a high
number of test can be computationally expensive.

Fuzz tests and especially fuzz targets need to be chosen carefully, as sometimes the fuzz test can fail due to
irrelevant inputs,

# Sources / Literature

go.dev

- [Go fuzz doc](https://go.dev/doc/security/fuzz/)
- [Go fuzz tutorial](https://go.dev/doc/tutorial/fuzz)
- [Source code fuzz.go](https://go.dev/src/internal/fuzz/fuzz.go)

blog articles

- [Best practices for go fuzzing in Go 1.18](https://faun.pub/best-practices-for-go-fuzzing-in-go-1-18-84eab46b70d8)
- [The state of Go Fuzzing: Did we already reach the peak](https://0x434b.dev/the-state-of-go-fuzzing-did-we-already-reach-the-peak/#native-go-fuzzing-is-it-advancing)
- [Finding bugs with go fuzzing](https://bitfieldconsulting.com/posts/bugs-fuzzing)

youtube videos

- [How to write a fuzz test | Demo](https://www.youtube.com/watch?v=y8Rpb3nrJn8&t=324s)
- [Introduction to Fuzzing](https://www.youtube.com/watch?v=-hc6LGA46Bg)
