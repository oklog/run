package run_test

import (
	"errors"
	"testing"
	"time"

	"github.com/oklog/run"
)

func TestZero(t *testing.T) {
	var g run.Group
	res := make(chan error)
	go func() { res <- g.Run() }()
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
	var g run.Group
	g.Add(func() error { return myError }, func(error) {})
	res := make(chan error)
	go func() { res <- g.Run() }()
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
	var g run.Group
	g.Add(func() error { return interrupt }, func(error) {})
	cancel := make(chan struct{})
	g.Add(func() error { <-cancel; return nil }, func(error) { close(cancel) })
	res := make(chan error)
	go func() { res <- g.Run() }()
	select {
	case err := <-res:
		if want, have := interrupt, err; want != have {
			t.Errorf("want %v, have %v", want, have)
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("timeout")
	}
}

func TestGroup_Run_AllActorsExecuteAndInterruptProperly(t *testing.T) {
	var (
		mock1Msg string
		mock2Msg string
		errMsg1  string
		errMsg2  string
	)

	var g run.Group

	// Actor 1: returns error
	g.Add(func() error {
		return errors.New("got error")
	}, func(err error) {
		errMsg1 = err.Error()
	})

	// Actor 2: completes successfully
	g.Add(func() error {
		return nil
	}, func(err error) {
		mock1Msg = "interrupted mock1"
	})

	// Actor 3: returns another error
	g.Add(func() error {
		return errors.New("got error 2")
	}, func(err error) {
		errMsg2 = err.Error()
	})

	// Actor 4: completes successfully
	g.Add(func() error {
		return nil
	}, func(err error) {
		mock2Msg = "interrupted mock2"
	})

	// Act
	err := g.Run()

	// Assert
	if err == nil {
		t.Fatal("expected non-nil error from Run, got nil")
	}

	if errMsg1 != "got error" && errMsg1 != "got error 2" {
		t.Errorf("expected errMsg1 to be 'got error', got '%s'", errMsg1)
	}

	if errMsg2 != "got error 2" && errMsg2 != "got error" {
		t.Errorf("expected errMsg2 to be 'got error 2', got '%s'", errMsg2)
	}

	if mock1Msg != "interrupted mock1" {
		t.Errorf("expected mock1Msg to be 'interrupted mock1', got '%s'", mock1Msg)
	}

	if mock2Msg != "interrupted mock2" {
		t.Errorf("expected mock2Msg to be 'interrupted mock2', got '%s'", mock2Msg)
	}
}
