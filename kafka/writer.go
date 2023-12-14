package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	Writer *kafka.Writer
	Params *[]TopicParams
}

func NewWriter(broker string, params []TopicParams, batchSize int) *Writer {
	writer := kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              batchSize,
		AllowAutoTopicCreation: true,
	}

	return &Writer{
		Writer: &writer,
		Params: &params,
	}
}

func (w *Writer) WriteMessages(ctx context.Context, payload []byte, topic string) error {
	return w.Writer.WriteMessages(ctx, kafka.Message{Value: payload, Topic: topic})
}
