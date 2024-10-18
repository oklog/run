package process

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/superblocksteam/run"
)

type process struct {
	signals []os.Signal

	run.ForwardCompatibility
}

func New(signals ...os.Signal) run.Runnable {
	if len(signals) == 0 {
		signals = append(signals, os.Interrupt, syscall.SIGTERM)
	}

	return &process{
		signals: signals,
	}
}

func (*process) Name() string { return "process manager" }
func (*process) Alive() bool  { return true }

func (p *process) Run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, p.signals...)
	defer cancel()

	<-ctx.Done()
	return nil
}
