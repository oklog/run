package run

import (
	"context"
	"log/slog"
)

// ForwardCompatibility provides a mechanism that allows Runnables to
// always be forwards compatible with future version of the  Runnable
// interface. The inspiration for this pattern comes from the Protobuf
// extension protoc-gen-grpc-go. We do not require it's embedding but
// it is highly recommended to ensure forwards compatibility.
type ForwardCompatibility struct{}

func (ForwardCompatibility) Run(context.Context) error   { panic("runnables must implement run") }
func (ForwardCompatibility) Close(context.Context) error { return nil }
func (ForwardCompatibility) Alive() bool                 { return true }
func (ForwardCompatibility) Name() string                { return "unknown" }
func (ForwardCompatibility) Fields() []slog.Attr         { return []slog.Attr{} }
