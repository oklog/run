package run

import (
	"context"

	"go.uber.org/zap"
)

type ForwardCompatibility struct{}

func (ForwardCompatibility) Run(context.Context) error { panic("runnables must implement run") }
func (ForwardCompatibility) Close()                    {}
func (ForwardCompatibility) Alive() bool               { return true }
func (ForwardCompatibility) Name() string              { return "unknown" }
func (ForwardCompatibility) Fields() []zap.Field       { return []zap.Field{} }
