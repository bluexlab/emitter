package emitter

import (
	"context"
	"sync"
	"time"
)

type Emitter struct {
	interval    time.Duration
	batchSize   int
	outbox      OutboxSource
	handler     Handler
	logger      Logger
	errorLogger Logger
	ctx         context.Context
	wg          *sync.WaitGroup
	stopSignal  chan any
}

func NewEmitter(ctx context.Context, outbox OutboxSource, handler HandlerFunc, opts ...Option) *Emitter {
	e := &Emitter{
		ctx:     ctx,
		outbox:  outbox,
		handler: handler,
	}

	for _, opt := range opts {
		opt(e)
	}
	if e.interval == 0 {
		e.interval = 30 * time.Second
	}
	if e.batchSize == 0 {
		e.batchSize = 10
	}

	return e
}

// Run starts the emitter.
func (e *Emitter) Run() {
	if e.wg != nil {
		return
	}

	e.withLogger(func(logger Logger) {
		logger.Printf("Emitter is running ...")
	})

	e.wg = &sync.WaitGroup{}
	e.wg.Add(1)
	e.stopSignal = make(chan any)
	go func() {
		defer e.wg.Done()

		ticker := time.NewTicker(e.interval)
		defer ticker.Stop()
		for {
			select {
			case <-e.stopSignal:
				return
			case <-ticker.C:
				e._Proc()
			}
		}
	}()
}

func (e *Emitter) _Proc() {
	for {
		select {
		case <-e.stopSignal:
			return
		default:
		}
		msgCount := e._ProcOneBatch(e.ctx)
		if msgCount <= 0 {
			return
		}
	}
}

func (e *Emitter) _ProcOneBatch(ctx context.Context) int {
	msgs, err := e.outbox.GetOutboxMsg(ctx, e.batchSize)
	if err != nil {
		e.withLogger(func(logger Logger) {
			logger.Printf("failed to GetOutboxMsg: %v", err)
		})
		return 0
	}
	if len(msgs) == 0 {
		return 0
	}

	recIDs, err := e.handler.Process(ctx, msgs...)
	if err != nil {
		e.withLogger(func(logger Logger) {
			logger.Printf("failed to process OutboxMsgs: %v", err)
		})
		return 0
	}

	err = e.outbox.DeleteOutboxMsg(ctx, recIDs...)
	if err != nil {
		e.withLogger(func(logger Logger) {
			logger.Printf("failed to DeleteOutboxMsg: %v", err)
		})
		return 0
	}

	return len(recIDs)
}

// Stop stops the emitter.
func (e *Emitter) Stop() {
	if nil == e.wg {
		return
	}

	close(e.stopSignal)
	e.wg.Wait()
	e.stopSignal = nil
	e.wg = nil
}

func (e *Emitter) withLogger(do func(Logger)) {
	if e.logger != nil {
		do(e.logger)
	}
}

func (e *Emitter) withErrorLogger(do func(Logger)) {
	if e.errorLogger != nil {
		do(e.errorLogger)
	} else {
		e.withLogger(do)
	}
}
