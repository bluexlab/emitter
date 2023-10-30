package emitter

import "context"

// OutboxSource provides the functionality to retrieve and deleted OutboxMsg.
type OutboxSource interface {
	GetOutboxMsg(ctx context.Context, batchSize int) ([]OutboxMsg, error)
	DeleteOutboxMsg(ctx context.Context, recIDs ...int64) error
}

type OutboxMsg struct {
	RecID int64
	Key   string
	Msg   []byte
}
