package ticker

import (
	"context"
	"sync"
	"time"

	"github.com/superblocksteam/run"
)

type runnable struct {
	fn       func()
	async    bool
	interval time.Duration
	wg       *sync.WaitGroup

	run.ForwardCompatibility
}

func New(interval time.Duration, fn func(), async bool) run.Runnable {
	return &runnable{
		fn:       fn,
		async:    async,
		interval: interval,
		wg:       &sync.WaitGroup{},
	}
}

func (r *runnable) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)

	for {
		select {
		case <-ticker.C:
			if r.async {
				r.wg.Add(1)
				go func() {
					defer r.wg.Done()
					r.fn()
				}()
			} else {
				r.fn()
			}
		case <-ctx.Done():
			ticker.Stop()
			r.wg.Wait()
			return nil
		}
	}
}

func (r *runnable) Name() string { return "ticker" }
