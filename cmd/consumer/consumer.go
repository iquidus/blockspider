package main

import (
	"context"
	"log"

	"github.com/iquidus/blockspider/kafka"
	gkafka "github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

func handleMessages(ctx context.Context, messages chan gkafka.Message, commits chan gkafka.Message) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message := <-messages:
			log.Printf("message fetched and sent to channel: %v \n", string(message.Value))
			select {
			case <-ctx.Done():
			case commits <- message:
			}
		}
	}
}

func main() {
	const (
		topic = "blocks"
	 	groupId = "example2"
		chanSize = 1000
	)

	var (
		ctx = context.Background()
		messages = make(chan gkafka.Message, chanSize)
		commits = make(chan gkafka.Message, chanSize)
	)

	kReader := kafka.NewReader([]string{"localhost:9092"}, topic, groupId)
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return kReader.FetchMessage(ctx, messages)
	})

	g.Go(func() error {
		return handleMessages(ctx, messages, commits)
	})

	g.Go(func() error {
		return kReader.CommitMessages(ctx, commits)
	})

	err := g.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}

