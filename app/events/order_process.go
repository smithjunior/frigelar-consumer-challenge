package event

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka-consumer-challenge/internal/config"
	infraHttp "kafka-consumer-challenge/internal/infrastructure/http"
	infraKafka "kafka-consumer-challenge/internal/infrastructure/kafka"
	"kafka-consumer-challenge/internal/usecase"
)

func Start(cfg *config.Config) {
	// Initialize Infrastructure (Adapters)
	// Kafka Producer
	producer, err := infraKafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// HTTP Client
	orderIntegration := infraHttp.NewHTTPOrderIntegration(cfg.Service.TargetServiceURL)

	// Initialize Use Case (Domain Service)
	orderProcessor := usecase.NewOrderProcessor(orderIntegration, producer, cfg.Kafka.TopicDLQ)

	// Initialize Kafka Consumer
	kafkaConsumer, err := infraKafka.NewConsumer(cfg.Kafka, orderProcessor)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer kafkaConsumer.Close()

	// Start Consumer in a goroutine
	go kafkaConsumer.Start()

	// Graceful Shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	log.Println("Received termination signal, shutting down...")
}
