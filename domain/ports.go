package domain

type EventProducer interface {
	Produce(topic string, key []byte, value []byte) error
}

type OrderIntegration interface {
	UpdateOrder(data []byte) error
}
