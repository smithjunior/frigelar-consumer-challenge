package kafka

import (
	"kafka-consumer-challenge/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": cfg.BootstrapServers,
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

	p, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(topic string, key []byte, value []byte) error {
	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

func (p *Producer) Close() {
	if p.producer != nil {
		p.producer.Close()
	}
}
