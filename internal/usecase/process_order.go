package usecase

import (
	"encoding/json"
	"log"

	"kafka-consumer-challenge/domain"
)

type OrderProcessor struct {
	orderIntegration domain.OrderIntegration
	producer         domain.EventProducer
	topicDLQ         string
}

func NewOrderProcessor(orderIntegration domain.OrderIntegration, producer domain.EventProducer, topicDLQ string) *OrderProcessor {
	return &OrderProcessor{
		orderIntegration: orderIntegration,
		producer:         producer,
		topicDLQ:         topicDLQ,
	}
}

func (p *OrderProcessor) Process(key []byte, data []byte) error {
	var payload domain.Order
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Failed to unmarshal payload: %v. Sending to DLQ...", err)
		p.sendToDLQ(nil, data)
		return nil 
	}

	log.Printf("Processing order: %s, Current Step: %s", payload.OrdemDeVenda, payload.EtapaAtual)

	if err := p.orderIntegration.UpdateOrder(data); err != nil {
		log.Printf("Error processing message: %v. Sending to DLQ...", err)
		p.sendToDLQ(key, data) 
		return nil 
	}

	log.Printf("Successfully processed and forwarded order: %s", payload.OrdemDeVenda)
	return nil
}

func (p *OrderProcessor) sendToDLQ(key []byte, value []byte) {
	if err := p.producer.Produce(p.topicDLQ, key, value); err != nil {
		log.Printf("Failed to produce to DLQ: %v", err)
	} else {
		log.Printf("Message sent to DLQ: %s", p.topicDLQ)
	}
}
