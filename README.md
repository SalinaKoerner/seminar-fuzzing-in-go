# Seminar: Fuzzing in Go

[Seminar Overview GitHub](https://github.com/sulzmann/Seminar/blob/main/winter24-25.md)

<!-- TOC -->
* [Seminar: Fuzzing in Go](#seminar-fuzzing-in-go)
  * [T3: Go's Built-In Fuzzer](#t3-gos-built-in-fuzzer)
* [Autonome Systeme Examples](#autonome-systeme-examples)
  * [Deadlock Example](#deadlock-example)
  * [Livelock Example](#livelock-example)
  * [Starvation Example](#starvation-example)
  * [Data Race Example](#data-race-example)
  * [Philo Example](#philo-example)
* [Conclusions](#conclusions)
    * [Questions:](#questions)
    * [Advatages of Go's Built-In Fuzzer](#advatages-of-gos-built-in-fuzzer)
    * [Limitations of Go's Built-In Fuzzer](#limitations-of-gos-built-in-fuzzer)
* [Sources / Literature](#sources--literature-)
<!-- TOC -->

## T3: Go's Built-In Fuzzer

Topic: How effective is Go fuzzing to detect concurrency bugs?

1. As a starting point, consider the bug scenarios and examples discussed in Autonome Systeme.

2. There are further fuzzing tools for concurrent Go. See below. Check out some of the examples used and.

3. Apply Go fuzzing and report your experiences.

---

# Autonome Systeme Examples
- Bugs scenarios are:
  - deadlock
  - livelock
  - starvation
  - data race

## Deadlock Example

```
package main

import "fmt"

func snd(ch chan int) {
var x int = 0
x++
ch <- x
}

func rcv(ch chan int) {
var x int
x = <-ch
fmt.Printf("received %d \n", x)

}

func main() {
var ch chan int = make(chan int)
go rcv(ch)   // R1
go snd(ch)   // S1
rcv(ch)      // R2

}
```

Fuzz test
```
package deadlock

import "testing"

func FuzzDeadlock(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, input int) {
		var ch chan int = make(chan int)

		go rcv(ch) // R1
		go snd(ch) // S1
		rcv(ch) //R2
	})
}
```
sample run
```
λ go test -fuzz=FuzzDeadlock
fuzz: elapsed: 0s, gathering baseline coverage: 0/1 completed
fuzz: elapsed: 0s, gathering baseline coverage: 1/1 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 7864 (2620/sec), new interesting: 8 (total: 9)
fuzz: elapsed: 6s, execs: 7864 (0/sec), new interesting: 8 (total: 9)
fuzz: elapsed: 9s, execs: 7864 (0/sec), new interesting: 8 (total: 9)
fuzz: elapsed: 10s, execs: 7994 (107/sec), new interesting: 8 (total: 9)
--- FAIL: FuzzDeadlock (10.22s)
    fuzzing process hung or terminated unexpectedly: exit status 2
    Failing input written to testdata\fuzz\FuzzDeadlock\0591344243e3314b
    To re-run:
    go test -run=FuzzDeadlock/0591344243e3314b
FAIL
exit status 1
FAIL    main/main/deadlock      10.503s
```

## Livelock Example

```package main

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
package livelock

import (
	"testing"
	"time"
)

func FuzzLivelock(f *testing.F) {
	f.Add(1)

	f.Fuzz(func(t *testing.T, input int) {
		go livelock()
		time.Sleep(1 * time.Second)
	})
}

```

Sample run

Does not detect a bug

## Starvation Example

```
package main

import "fmt"
import "time"

  func snd(ch chan int) {
    var x int = 0
    for {
      x++
      ch <- x
      time.Sleep(1 * 1e9)
  }

}

  func rcv(ch chan int) {
  var x int
  for {
    x = <-ch
    fmt.Printf("received %d \n", x)
  }

}

func main() {
var ch chan int = make(chan int)
go rcv(ch)   // R1
go snd(ch)   // S1
rcv(ch)      // R2

}
```

Fuzz test
```
package starvation

import (
	"testing"
)

func FuzzStarvation(f *testing.F) {
	f.Add(0)

	f.Fuzz(func(t *testing.T, _ int) {
		var ch chan int = make(chan int)
		
		go rcv(ch) // R1
		go snd(ch) // S1
		rcv(ch)    // R2

	})
}

```

Sample run
```
λ go test -fuzz=FuzzStarvation
fuzz: elapsed: 0s, gathering baseline coverage: 0/5 completed
fuzz: elapsed: 3s, gathering baseline coverage: 0/5 completed
fuzz: elapsed: 6s, gathering baseline coverage: 0/5 completed
fuzz: elapsed: 9s, gathering baseline coverage: 0/5 completed
failure while testing seed corpus entry: FuzzStarvation/seed#0
fuzz: elapsed: 10s, gathering baseline coverage: 0/5 completed
--- FAIL: FuzzStarvation (10.07s)
    fuzzing process hung or terminated unexpectedly: exit status 2
FAIL
exit status 1
FAIL    main/main/starvation    10.352s
```

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

Sample run
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

## Philo Example

The example: The Dining Philosophers Problem

N philosophers sit at a table with a total of N forks. Each philosopher requires 2 forks to eat.

Possible bugs:
- deadlock 
- starvation

```
package main

import "fmt"
import "time"

func philo(id int, forks chan int) {

	for {
		<-forks
		<-forks
		fmt.Printf("%d eats \n", id)
		time.Sleep(1 * 1e9)
		forks <- 1
		forks <- 1

		time.Sleep(1 * 1e9) // think

	}

}

func main() {
  var forks = make(chan int, 3)
  forks <- 1
  forks <- 1
  forks <- 1
  go philo(1, forks)
  go philo(2, forks)
  philo(3, forks)
}
```
The fuzz test:

```
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
```

The sample output of the fuzz test:

```
λ go test -fuzz=FuzzPhilo

fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 249422 (82902/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 6s, execs: 430488 (60524/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 9s, execs: 628640 (65775/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 12s, execs: 817335 (62933/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 15s, execs: 1072511 (84920/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 18s, execs: 1304837 (77489/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 21s, execs: 1532728 (76017/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 24s, execs: 1729492 (65671/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 27s, execs: 1821319 (30686/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 30s, execs: 1935366 (37577/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 33s, execs: 2094398 (53638/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 36s, execs: 2195395 (33628/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 39s, execs: 2315483 (40008/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 42s, execs: 2437627 (40782/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 45s, execs: 2557238 (39877/sec), new interesting: 14 (total: 17)
fuzz: elapsed: 47s, execs: 2618761 (32256/sec), new interesting: 14 (total: 17)
--- FAIL: FuzzPhilo (46.91s)
fuzzing process hung or terminated unexpectedly: exit status 2
Failing input written to testdata\fuzz\FuzzPhilo\ace36d9332763359
To re-run:
go test -run=FuzzPhilo/ace36d9332763359
FAIL
exit status 1
FAIL    main/main/philo 47.207s
```

# Conclusions

### Questions:
- What is it good at detecting, where does it fail to detect bugs?
- What are the limitations of the built-in fuzzer?

# Sources / Literature 
- [go.dev - Tutorial: Getting started with fuzzing ](https://go.dev/doc/tutorial/fuzz)