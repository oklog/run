package run

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestZero(t *testing.T) {
	var g group
	res := make(chan error)
	go func() { res <- g.run() }()
	select {
	case err := <-res:
		if err != nil {
			t.Errorf("%v", err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout")
	}
}

func TestOne(t *testing.T) {
	myError := errors.New("foobar")
	var g group
	g.add(func(context.Context) error { return myError }, func() {})
	res := make(chan error)
	go func() { res <- g.run() }()
	select {
	case err := <-res:
		if want, have := myError, err; want != have {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout")
	}
}

func TestMany(t *testing.T) {
	interrupt := errors.New("interrupt")
	var g group
	g.add(func(context.Context) error { return interrupt }, func() {})
	cancel := make(chan struct{})
	g.add(func(context.Context) error { <-cancel; return nil }, func() { close(cancel) })
	res := make(chan error)
	go func() { res <- g.run() }()
	select {
	case err := <-res:
		if want, have := interrupt, err; want != have {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("timeout")
	}
}
