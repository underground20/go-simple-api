package kafka

import (
	"app/lib/logger"
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler func(kafka.Message)
	logger  *slog.Logger
	wg      *sync.WaitGroup
}

func NewConsumer(
	handler func(message kafka.Message),
	logger *slog.Logger,
	brokers []string,
	topic string,
	groupId string,
) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:           brokers,
			Topic:             topic,
			GroupID:           groupId,
			StartOffset:       kafka.FirstOffset,
			CommitInterval:    1 * time.Second,
			HeartbeatInterval: 3 * time.Second,
			SessionTimeout:    30 * time.Second,
		}),
		handler: handler,
		logger:  logger,
		wg:      &sync.WaitGroup{},
	}
}

func (c *Consumer) ReadMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Error("failed to read message:", logger.Err(err))
				continue
			}
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				c.handler(m)
			}()
		}
	}
}

func (c *Consumer) Close() error {
	err := c.reader.Close()
	c.wg.Wait()
	return err
}
