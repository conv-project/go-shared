package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/conv-project/conversion-service/pkg/logger"
)

// Producer represents Kafka producer.
type Producer struct {
	producer sarama.SyncProducer
	brokers  []string
}

// NewProducer creates a new Kafka producer.
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	logger.Info("Connected to Kafka producer", zap.Strings("brokers", brokers))

	return &Producer{
		producer: producer,
		brokers:  brokers,
	}, nil
}

// SendMessage sends a message to Kafka topic.
func (p *Producer) SendMessage(topic string, key string, value []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	logger.Debug("Message sent to Kafka",
		zap.String("topic", topic),
		zap.String("key", key),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
	)

	return nil
}

// Close closes Kafka producer.
func (p *Producer) Close() error {
	if p.producer != nil {
		if err := p.producer.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka producer: %w", err)
		}
		logger.Info("Kafka producer closed")
	}
	return nil
}

// Consumer represents Kafka consumer.
type Consumer struct {
	consumer sarama.ConsumerGroup
	brokers  []string
	group    string
	topics   []string
	handler  sarama.ConsumerGroupHandler
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(brokers []string, groupID string, topics []string, handler sarama.ConsumerGroupHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	logger.Info("Connected to Kafka consumer",
		zap.Strings("brokers", brokers),
		zap.String("group", groupID),
		zap.Strings("topics", topics),
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer: consumer,
		brokers:  brokers,
		group:    groupID,
		topics:   topics,
		handler:  handler,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start starts consuming messages from Kafka.
func (c *Consumer) Start() {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			if err := c.consumer.Consume(c.ctx, c.topics, c.handler); err != nil {
				logger.Error("Error from consumer", zap.Error(err))
			}

			if c.ctx.Err() != nil {
				return
			}
		}
	}()

	logger.Info("Kafka consumer started",
		zap.String("group", c.group),
		zap.Strings("topics", c.topics),
	)
}

// Close stops consuming messages and closes Kafka consumer.
func (c *Consumer) Close() error {
	c.cancel()
	c.wg.Wait()

	if c.consumer != nil {
		if err := c.consumer.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka consumer: %w", err)
		}
		logger.Info("Kafka consumer closed")
	}
	return nil
}
