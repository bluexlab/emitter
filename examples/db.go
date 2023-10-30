package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bluexlab/emitter"
)

type DB struct {
	data map[int]emitter.OutboxMsg
}

func NewDB(size int) *DB {
	data := make(map[int]emitter.OutboxMsg, size)
	for i := 1; i <= size; i++ {
		data[i] = emitter.OutboxMsg{
			RecID: int64(i),
			Key:   fmt.Sprintf("key-%d", i),
			Msg:   []byte(fmt.Sprintf("msg-%d", i)),
		}
	}

	return &DB{data: data}
}

func (d *DB) GetOutboxMsg(ctx context.Context, batchSize int) ([]emitter.OutboxMsg, error) {
	var res []emitter.OutboxMsg
	for _, v := range d.data {
		res = append(res, v)
		batchSize--
		if batchSize == 0 {
			break
		}
	}

	return res, nil
}

func (d *DB) DeleteOutboxMsg(ctx context.Context, recIDs ...int64) error {
	for _, recID := range recIDs {
		slog.Info(fmt.Sprintf("Deleting record %d", recID))
		delete(d.data, int(recID))
	}

	return nil
}
