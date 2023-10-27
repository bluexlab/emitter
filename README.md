# Emitter

Emitter is the implementation of the Outbox Pattern for handling outbox message from database.

## Installation

```bash
go get -u github.com/bluexlab/emitter
```

## Default Usage

```go
// Kafka Writer
kafkaWriter := emitter.NewKafkaWriter([]string{"broker:9092"}, "topic", 10)
defer func() { _ = kafkaWriter.Close() }()

ctx, cancel := context.WithCancel(context.Background())
notificationEmitter := emitter.NewEmitter(ctx, outboxSourceImpl, emitter.KafkaHandler(kafkaWriter))
defer notificationEmitter.Stop()

notificationEmitter.Run()
```

## Message Handler
Emitter leverage `Handler` to process the messages.
User can provide custom handler as long as it follows `HandlerFunc` function signature:

```go   
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
```

## Options

By default, Emitter polls OutboxMsg every **30** second.
User can provide Options to change the default behavior:

```go
notificationEmitter := emitter.NewOutboxEmitter(ctx, source, handler,
    emitter.WithInterval(10 * time.Second),
    emitter.WithBatchSize(5),
)
```
