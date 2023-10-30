package emitter

import (
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

// NewKafkaWriter creates a new kafka writer.
func NewKafkaWriter(brokers []string, topic string, batchSize int) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		BatchSize:    batchSize,
	}
}

// KafkaHandler writes messages to kafka with a kafka writer.
func KafkaHandler(writer *kafka.Writer) HandlerFunc {
	return func(ctx context.Context, msgs ...OutboxMsg) ([]int64, error) {
		kafkaMsgs := make([]kafka.Message, len(msgs))
		recIDs := make([]int64, len(msgs))
		for i := range msgs {
			kafkaMsgs[i].Key = []byte(msgs[i].Key)
			kafkaMsgs[i].Value = msgs[i].Msg
			recIDs[i] = msgs[i].RecID
		}

		err := writer.WriteMessages(ctx, kafkaMsgs...)
		if err != nil {
			return nil, err
		}

		return recIDs, nil
	}
}
