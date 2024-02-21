package run

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

type enriched struct {
	ForwardCompatibility
}

func (*enriched) Name() string                { return "runnable" }
func (*enriched) Run(context.Context) error   { return nil }
func (*enriched) Alive() bool                 { return true }
func (*enriched) Close(context.Context) error { return nil }
func (*enriched) Fields() []slog.Attr         { return []slog.Attr{slog.String("foo", "bar")} }

type basic struct {
	ForwardCompatibility
}

func (*basic) Run(context.Context) error   { return nil }
func (*basic) Alive() bool                 { return false }
func (*basic) Close(context.Context) error { return nil }

func TestFields(t *testing.T) {
	assert.Equal(t, []slog.Attr{slog.String("foo", "bar")}, (&enriched{}).Fields())
	assert.Equal(t, []slog.Attr{}, (&basic{}).Fields())
}

func TestName(t *testing.T) {
	assert.Equal(t, "runnable", (&enriched{}).Name())
	assert.Equal(t, "unknown", (&basic{}).Name())
}

func TestAlive(t *testing.T) {
	g := New()

	g.Always(&enriched{})
	assert.True(t, g.Alive())

	g.Add(false, &basic{})
	assert.True(t, g.Alive())

	g.Add(true, &basic{})
	assert.False(t, g.Alive())
}

func TestAdd(t *testing.T) {
	g := New()

	g.Add(true, &enriched{})
	g.Add(false, &basic{})

	assert.Len(t, g.runnables, 1)
}
