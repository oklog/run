package waitgroup

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitAlive(t *testing.T) {
	assert.True(t, NewWait(&sync.WaitGroup{}).Alive())
}

func TestWait(t *testing.T) {
	wg := sync.WaitGroup{}
	done := make(chan struct{})
	runnable := NewWait(&wg)

	go func() {
		runnable.Run(context.Background())
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
	}()

	time.Sleep(50 * time.Millisecond) // Ensure goroutine has time to execute
	runnable.Close(nil)
	<-done
}
