package preempt

import (
	"context"
	"log/slog"
	"time"

	"github.com/superblocksteam/run"
)

type preempt struct {
	after time.Duration

	run.ForwardCompatibility
}

func New(after time.Duration) run.Runnable {
	return &preempt{
		after: after,
	}
}

func (*preempt) Name() string { return "preempter" }
func (*preempt) Alive() bool  { return true }

func (p *preempt) Run(ctx context.Context) error {
	select {
	case <-time.After(p.after):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *preempt) Fields() []slog.Attr {
	return []slog.Attr{
		slog.Duration("after", p.after),
	}
}
