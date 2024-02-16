# run

[![Go Reference](https://pkg.go.dev/badge/github.com/superblocksteam/run.svg)](https://pkg.go.dev/github.com/superblocksteam/run) [![Go Report Card](https://goreportcard.com/badge/github.com/superblocksteam/run)](https://goreportcard.com/report/github.com/superblocksteam/run)

## Motivation

A program will likely have many long lived daemons running concurrently in go routines. When the program needs to shutdown, it should do so gracefully by telling and waiting for every concurrent daemon to shutdown. There are many ways in which these semantics could be implemented and [Peter Bourgon](https://github.com/peterbourgon) compares and contrasts many of them in [this](https://www.youtube.com/watch?v=LHe1Cb_Ud_M) video. He recommends the pattern implemented in [`oklog/run`](https://github.com/oklog/run). After 5+ years of using this project, there were various enhancements and wrappers that we made over time to improve the developer experience. This project open sources them. You could think of this as the next generation of `oklog/run`.

## Usage

Every daemon must implement the `Runnable` interface. A set of Runnable types are then registered to a group and executed. We demonstrate this with a runnable whose job it is to sleep for a specified duration and then exit.

```go
package main

import (
  "time"
  "context"

  "github.com/superblocksteam/run"
)

type preempt struct {
  after time.Duration

  run.ForwardCompatibility
}

func (p *preempt) Run(ctx context.Context) error {
  select {
    case <-time.After(p.after):
    return nil
  case <-ctx.Done():
    return ctx.Err()
  }
}

func main() {
  run.Always(&preempt{5 * time.Second})
  run.Run()
}
```

## Contributions

The `contrib` package provides runnable implementations for common use cases.

## Comparisons

The main comparison to draw is with `oklog/run`. This package is different in the following ways 👇

1. Do not invoke an actor's interrupt synchronously. Rather, invoke it concurrently and synchronize after so that with the help of logs, we can detect which actors may be blocking and/or implemented incorrectly. This does not modify the semantics but aids in debugging.
2. Rather than register an actor as two anonymous functions, move it to an interface named `Runnable` that types can implement.
3. Add an optional embeddable struct to allow for forward compatibility of the Runnable interface in the same way that `protoc-gen-grpc-go` does.
4. Allow actors to be added conditionally based on some boolean. This is useful when actors are conditionally executed based on the result of a runtime parameter.
5. Add additional methods, other than `Run` and `Close` to the `Runnable` interface to improve the functionality. See the documentation for more details.
