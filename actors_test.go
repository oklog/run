package run

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestContextHandler(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var rg Group
	rg.Add(ContextHandler(ctx))
	errc := make(chan error, 1)
	go func() { errc <- rg.Run() }()
	cancel()
	select {
	case err := <-errc:
		if want, have := context.Canceled, err; !errors.Is(have, want) {
			t.Errorf("error: want %v, have %v", want, have)
		}
	case <-time.After(time.Second):
		t.Errorf("timeout waiting for error after cancel")
	}
}

func TestSignalError(t *testing.T) {
	testc := make(chan os.Signal, 1)
	ctx := putTestSigChan(context.Background(), testc)

	var rg Group
	rg.Add(SignalHandler(ctx, os.Interrupt))
	testc <- os.Interrupt
	err := rg.Run()

	var sigerr *SignalError
	if want, have := true, errors.As(err, &sigerr); want != have {
		t.Errorf("errors.As(err, &sigerr): want %v, have %v", want, have)
	}

	if sigerr != nil {
		if want, have := os.Interrupt, sigerr.Signal; want != have {
			t.Errorf("sigerr.Signal: want %v, have %v", want, have)
		}
	}

	if sigerr := &(SignalError{}); !errors.As(err, &sigerr) {
		t.Errorf("errors.As(err, <inline sigerr>): failed")
	}

	if want, have := true, errors.As(err, &(SignalError{})); want != have {
		t.Errorf("errors.As(err, &(SignalError{})): want %v, have %v", want, have)
	}

	if want, have := true, errors.Is(err, ErrSignal); want != have {
		t.Errorf("errors.Is(err, ErrSignal): want %v, have %v", want, have)
	}
}

func TestSignalHandlerNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var rg Group
	rg.Add(SignalHandler(ctx, os.Interrupt))
	cancel()
	t.Logf("%v", rg.Run())
}
