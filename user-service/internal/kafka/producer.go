package kafka

import (
	"context"
	"encoding/json"
	"time"

	"user-service/internal/config"

	"github.com/segmentio/kafka-go"
)

type Producer interface {
	Publish(ctx context.Context, topic string, key string, value any) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *config.Config) Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.KafkaProducer.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Publish(ctx context.Context, topic string, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
