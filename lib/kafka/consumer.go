package kafka

import (
	"app/lib/logger"
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler func(kafka.Message) error
	logger  *slog.Logger
	wg      *sync.WaitGroup
}

func NewConsumer(
	handler func(message kafka.Message) error,
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
			HeartbeatInterval: 3 * time.Second,
			SessionTimeout:    30 * time.Second,
			Logger:            kafka.LoggerFunc(logger.Info),
			ErrorLogger:       kafka.LoggerFunc(logger.Error),
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
			c.logger.Info("Shutting down consumer gracefully...")
			c.wg.Wait()
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					continue
				}
				c.logger.Error("failed to read message:", logger.Err(err))
				continue
			}

			c.wg.Add(1)
			go func(m kafka.Message) {
				defer c.wg.Done()
				err := c.handler(m)
				if err != nil {
					c.logger.Error("failed to handle message:", logger.Err(err))
					return
				}

				commitErr := c.reader.CommitMessages(ctx, m)
				if commitErr != nil {
					c.logger.Error("Failed to commit offset",
						slog.String("topic", m.Topic),
						slog.Int("partition", m.Partition),
						slog.Int64("offset", m.Offset),
						logger.Err(commitErr))
				} else {
					c.logger.Debug("Committed offset", slog.Int64("offset", m.Offset))
				}
			}(msg)
		}
	}
}

func (c *Consumer) Close() error {
	err := c.reader.Close()
	c.wg.Wait()
	return err
}
