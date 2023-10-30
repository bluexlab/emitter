package emitter

import (
	"time"
)

// Option is a function that configures an Emitter.
type Option func(emitter *Emitter)

// WithInterval sets the interval for the emitter.
func WithInterval(interval time.Duration) Option {
	return func(e *Emitter) {
		e.interval = interval
	}
}

// WithBatchSize sets the batch size for the emitter.
func WithBatchSize(size int) Option {
	return func(e *Emitter) {
		e.batchSize = size
	}
}

// WithLogger sets the logger for the antenna.
func WithLogger(logger Logger) Option {
	return func(e *Emitter) {
		e.logger = logger
	}
}

// WithErrorLogger sets the error logger for the antenna.
func WithErrorLogger(logger Logger) Option {
	return func(e *Emitter) {
		e.errorLogger = logger
	}
}
