package emitter

import (
	"context"
)

// Handler process the OutboxMsg retrieve from OutboxSource.
type Handler interface {
	Process(ctx context.Context, msgs ...OutboxMsg) ([]int64, error)
}

// HandlerFunc is an adapter to allow the use of
// ordinary functions as a OutboxMsg handler. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(ctx context.Context, msgs ...OutboxMsg) ([]int64, error)

// Process calls f(ctx, msg).
func (f HandlerFunc) Process(ctx context.Context, msgs ...OutboxMsg) ([]int64, error) {
	return f(ctx, msgs...)
}
