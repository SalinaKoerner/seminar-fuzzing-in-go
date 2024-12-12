## How the Go fuzzer works

1. start with seed/given inputs
2. mutate inputs
3. monitors program for errors
4. write failing inputs into file

## Structure of a fuzz test

- name: Fuzz{FunctionToFuzz}
- file: {FunctionToFuzz}_test.go

![fuzzTestStructure.png](fuzzTestStructure.png)

- **Fuzz test**: entire test with all its required parts, including fuzz target, fuzzing arguments
  and seed additions.

- **Fuzz target**: function being executed with the corpus entries and the generated inputs.

- **Fuzzing arguments**: data types that are being passed to the fuzzing function.

    - are mutated to generate new inputs.

- **Seed corpus addition**: arguments provided to be added to the corpus.

## Output of a fuzz test

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

### Meaning of the output

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

# Pros and Cons / Advantages and Limitations

| Advantages                                                                              | Disadvantages                                                 |
|-----------------------------------------------------------------------------------------|---------------------------------------------------------------|
| can find bugs which unit tests can't find                                               | not a lot of documentation available, because it's fairly new |
| seems to be good at finding data races which the data race detector doesn't always find | doesn't always find every bug                                 |
| can be useful for parsing errors or incorrect rune decoding                             | only few data types are supported as fuzzing arguments        |
| can find bugs the programmer didn't even think about or didn't write a test for         | can take a long time and a lot of resources                   |