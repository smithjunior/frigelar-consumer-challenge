package kafka

import (
	"log"
	"time"

	"kafka-consumer-challenge/internal/config"
	"kafka-consumer-challenge/internal/usecase"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
	consumer  *kafka.Consumer
	processor *usecase.OrderProcessor
	cfg       config.KafkaConfig
	running   bool
}

func NewConsumer(cfg config.KafkaConfig, processor *usecase.OrderProcessor) (*Consumer, error) {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
		"group.id":          cfg.GroupId,
		"auto.offset.reset": "earliest",
	}

	// Setup SASL if configured
	if cfg.SecurityProtocol != "" {
		_ = kafkaConfig.SetKey("security.protocol", cfg.SecurityProtocol)
	}
	if cfg.SaslMechanism != "" {
		_ = kafkaConfig.SetKey("sasl.mechanism", cfg.SaslMechanism)
		_ = kafkaConfig.SetKey("sasl.username", cfg.Username)
		_ = kafkaConfig.SetKey("sasl.password", cfg.Password)
	}

	c, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{cfg.Topic}, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:  c,
		processor: processor,
		cfg:       cfg,
		running:   true,
	}, nil
}

func (c *Consumer) Start() {
	log.Printf("Starting consumer for topic: %s", c.cfg.Topic)
	for c.running {
		msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
		if err == nil {
			log.Printf("Received message on %s: %s\n", msg.TopicPartition, string(msg.Value))

			if err := c.processor.Process(msg.Key, msg.Value); err != nil {
				log.Printf("Error processing message (unhandled): %v", err)
			}
		} else {
			if !err.(kafka.Error).IsFatal() {
				continue
			}
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}

func (c *Consumer) Close() {
	c.running = false
	if c.consumer != nil {
		c.consumer.Close()
	}
}
