# run

[![Go Reference](https://pkg.go.dev/badge/github.com/superblocksteam/run.svg)](https://pkg.go.dev/github.com/superblocksteam/run) [![Go Report Card](https://goreportcard.com/badge/github.com/superblocksteam/run)](https://goreportcard.com/report/github.com/superblocksteam/run)

## Motivation

A program will likely have many long lived daemons running concurrently in go routines. An HTTP or gRPC server, Kafka consumer, etc. When the program needs to shutdown, it should do so gracefully by telling and waiting for every concurrent daemon to shutdown. There are many ways in which these semantics could be implemented and [Peter Bourgon](https://github.com/peterbourgon) compares and contrasts many of them in [this](https://www.youtube.com/watch?v=LHe1Cb_Ud_M) video. He recommends the pattern implemented in [`oklog/run`](https://github.com/oklog/run). After 5+ years of using this project, there were various enhancements and wrappers that we made over time to improve the developer experience. This project open sources them. You could think of this as the next generation of `oklog/run`.

## Semantics

1. Invoke `Run` for every runnable in the group concurrently.
2. Wait for the first one to return.
3. Cancel the context which will signal implementations not using `Close` to shutdown.
4. Invoke `Close` for every runnable in the group concurrently. If this method is not implemented, it means that `ForwardCompatibility`, whose implementation immediately returns `nil`, is embedded.
5. Wait for every `Close` to return
6. Wait for every `Run` to return
7. Return the error, if any, from step 2.

## Usage

Note that while we use the term _daemon_, this also includes clients that must be properly shutdown down. Examples of these include a syslogger or OTEL exporter that must be properly flushed, a Redis or PostgreSQL connection that must be closed, etc.

Every daemon must implement the [`Runnable`](https://pkg.go.dev/github.com/superblocksteam/run#Runnable) interface. While there are many methods on this interface, the only required method is `Run(context.Context)`. This is made possible through the embeddable [`ForwardCompatibility`](https://pkg.go.dev/github.com/superblocksteam/run#ForwardCompatibility). Here is a example implementing a process manager.

```go
type process struct {
  run.ForwardCompatibility
}

func (*process) Run(ctx context.Context) error {
  ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
  defer cancel()

  <-ctx.Done()
  return ctx.Err()
}
```

In this example, we run until (1) we need to exit because we received a signal or (2) we're being told to exit by the parent. There are some examples where it's cleaner to make use of the `Close(context.Context)` method on the `Runnable` interface. An HTTP server is such an example.

```go
type server struct {
  *http.Server
  run.ForwardCompatibility
}

func (s *server) Run(context.Context) error {
  return s.ListenAndServe() 
}

func (s *server) Close(ctx context.Context) error {
  return s.Shutdown(ctx)
}
```

A set of Runnables are then run as a group.

```go
group := New()

group.Always(new(process))
group.Always(new(server))

group.Run()
```

## Contributions

The `contrib` package provides runnable implementations for common use cases.

## Comparisons

The main comparison to draw is with `oklog/run`. While this project build on the original codebase, it is different in the following ways (terms are relative to `oklog/run`) 👇

1. Actor's interrupts are not invoked synchronously. Rather, the code invokes them concurrently and then waits for all of them to return. This allows all actors to make progress even if there is one that is "stuck" gracefully shutting down. This technical modifies the semantics if clients relied on an ordered shutdown. However, in our many years of using this, we never relied on that.
2. Provide the `Runnable` interface as syntactic sugar to make it easier to implement and wire up.
3. If a client wants to be more idiomatic, pass a context into an actor's execution allowing for graceful shutdown without implementing an interrupt.
4. Introduce an embeddable type allowing for forward compatibility of the Runnable interface in the same way that `protoc-gen-grpc-go` does. There is no way to require its embedding however.
5. Allow actors to be added. This is useful when actors are conditionally executed based on the result of a runtime parameter (i.e. `viper.GetBool("my_http_server.enabled")`).
6. Add new methods, in addition to `Run` and `Close`, to the `Runnable` interface to improve observability and assess health. See the [documentation](https://pkg.go.dev/github.com/superblocksteam/run#Runnable) for more details.
7. Allow the shutdown to be preempted by a context.
