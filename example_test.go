package run_test

import (
	"time"

	"github.com/superblocksteam/run"
	"github.com/superblocksteam/run/contrib/preempt"
)

func Example() {
	run.Add(true, preempt.New(100*time.Millisecond))

	run.Run()
}
