package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bluexlab/emitter"
)

func main() {
	slog.Info("Service started.")

	// Create context that listens for the interrupt signal from the OS.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	msgHandler := func(ctx context.Context, msgs ...emitter.OutboxMsg) ([]int64, error) {
		recIDs := make([]int64, len(msgs))
		for i := range msgs {
			recIDs[i] = msgs[i].RecID
			slog.Info(fmt.Sprintf("Processing OutboxMsg, Key: %s, Message: %s", msgs[i].Key, string(msgs[i].Msg)))
		}

		return recIDs, nil
	}

	outboxEmitter := emitter.NewEmitter(ctx, NewDB(23), msgHandler,
		emitter.WithBatchSize(5),
		emitter.WithInterval(5*time.Second),
		emitter.WithLogger(emitter.LoggerFunc(logInfo)),
		emitter.WithErrorLogger(emitter.LoggerFunc(logError)))
	defer outboxEmitter.Stop()

	outboxEmitter.Run()

	<-ctx.Done()
	cancel()

	slog.Info("Service stopped.")
}

func logInfo(msg string, args ...any) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func logError(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
}
