package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafka.Reader
}

func NewReader(brokers []string, topic string, groupId string) *Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
	})

	return &Reader{
		Reader: reader,
	}
}

func (k *Reader) FetchMessage(ctx context.Context, messages chan<- kafka.Message) error {
	for {
		message, err := k.Reader.FetchMessage(ctx)
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case messages <- message:
			log.Printf("message fetched and sent to channel: %v \n", string(message.Value))
		}
	}
}

func (k *Reader) CommitMessages(ctx context.Context, commits <-chan kafka.Message) error {
	for {
		select {
		case <-ctx.Done():
		case message := <-commits:
			err := k.Reader.CommitMessages(ctx, message)
			if err != nil {
				return err
			}
			log.Printf("committed a msg: %v \n", string(message.Value))
		}
	}
}
