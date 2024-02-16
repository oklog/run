package run_test

import (
	"time"

	"github.com/superblocksteam/run"
	"github.com/superblocksteam/run/runnables"
)

func Example_run() {
	run.Add(true, runnables.Process())
	run.Add(true, runnables.Preempt(5*time.Second))

	if err := run.Run(); err != nil {
		panic(err)
	}

	// Output:
	// {"time":"2024-02-16T12:52:22.27183-05:00","level":"INFO","msg":"started","name":"preempter","after":5000000000,"method":"run"}
	// {"time":"2024-02-16T12:52:22.271831-05:00","level":"INFO","msg":"started","name":"process manager","method":"run"}
	// {"time":"2024-02-16T12:52:27.273057-05:00","level":"INFO","msg":"returned","name":"preempter","after":5000000000,"method":"run"}
	// {"time":"2024-02-16T12:52:27.273327-05:00","level":"INFO","msg":"started","name":"preempter","after":5000000000,"method":"close"}
	// {"time":"2024-02-16T12:52:27.273381-05:00","level":"INFO","msg":"returned","name":"preempter","after":5000000000,"method":"close"}
	// {"time":"2024-02-16T12:52:27.273459-05:00","level":"INFO","msg":"started","name":"process manager","method":"close"}
	// {"time":"2024-02-16T12:52:27.27352-05:00","level":"INFO","msg":"returned","name":"process manager","method":"close"}
	// {"time":"2024-02-16T12:52:27.275005-05:00","level":"INFO","msg":"returned","name":"process manager","method":"run"}
}
