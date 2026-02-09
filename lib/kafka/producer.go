package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Message interface {
	Key() []byte
	Value() []byte
}

type Producer interface {
	SendMessage(ctx context.Context, msg Message) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
			Async:    true,
		},
	}
}

func (p *KafkaProducer) SendMessage(ctx context.Context, msg Message) error {
	err := p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   msg.Key(),
			Value: msg.Value(),
		},
	)
	if err != nil {
		log.Println("failed to write messages:", err)
		return err
	}
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
