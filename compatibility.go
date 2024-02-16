package run

import "go.uber.org/zap"

type ForwardCompatibility struct{}

func (ForwardCompatibility) Run() error {
	panic("all runnables must implement the run method")
}

func (ForwardCompatibility) Close(error) error {
	panic("all runnables must implement the close method")
}

func (ForwardCompatibility) Alive() bool         { return true }
func (ForwardCompatibility) Name() string        { return "unknown" }
func (ForwardCompatibility) Fields() []zap.Field { return []zap.Field{} }
func (ForwardCompatibility) compatibility()      {}
