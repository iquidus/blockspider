package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	Writer *kafka.Writer
}

func NewWriter(broker string, topic *string, batchSize int) *Writer {
	var writer kafka.Writer
	if topic != nil {
		writer = kafka.Writer{
			Addr: kafka.TCP(broker),
			Topic: *topic,
			Balancer: &kafka.LeastBytes{},
			BatchSize: batchSize,
		}
	} else {
		writer = kafka.Writer{
			Addr: kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
			BatchSize: batchSize,
		}
	}
	
	return &Writer{
		Writer: &writer,
	}
}

func (w *Writer) WriteMessages(ctx context.Context, payload []byte) error {
	return w.Writer.WriteMessages(ctx, kafka.Message{Value: payload})
}

func (w *Writer) WriteMessagesWithTopic(ctx context.Context, payload []byte, topic string) error {
	return w.Writer.WriteMessages(ctx, kafka.Message{Value: payload, Topic: topic})
}