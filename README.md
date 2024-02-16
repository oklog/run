# run

[![Go Reference](https://pkg.go.dev/badge/github.com/superblocksteam/run.svg)](https://pkg.go.dev/github.com/superblocksteam/run)

After 5+ years of using the amazing [`oklog/run`](https://github.com/oklog/run) package that [Peter Bourgon](https://github.com/peterbourgon) gave us, there were various enhancements and wrappers that we made over time to improve the developer experience. This project open sources these. Specifically, here are the differences with [`oklog/run`](https://github.com/oklog/run) 👇

1. Do not invoke an actor's interrupt synchronously. Rather, invoke it concurrently and synchronize after so that with the help of logs, we can detect which actors may be blocking and/or implemented incorrectly. This does not modify the semantics but aids in debugging.
2. Rather than register an actor as two anonymous functions, move it to an interface named `Runnable` that types can implement.
3. Add an optional embeddable struct to allow for forward compatibility of the Runnable interface in the same way that `protoc-gen-grpc-go` does.
4. Allow actors to be added conditionally based on some boolean. This is useful when actors are conditionally executed based on the result of a runtime parameter.
5. Add additional methods, other than `Run` and `Close` to the `Runnable` interface to improve the functionality. See the documentation for more details.
