package usecase

import (
	"errors"
	"testing"
)

// -- Mocks --

type MockEventProducer struct {
	ProduceFunc func(topic string, key []byte, value []byte) error
}

func (m *MockEventProducer) Produce(topic string, key []byte, value []byte) error {
	if m.ProduceFunc != nil {
		return m.ProduceFunc(topic, key, value)
	}
	return nil
}

type MockOrderIntegration struct {
	UpdateOrderFunc func(data []byte) error
}

func (m *MockOrderIntegration) UpdateOrder(data []byte) error {
	if m.UpdateOrderFunc != nil {
		return m.UpdateOrderFunc(data)
	}
	return nil
}

func TestOrderProcessor_Process(t *testing.T) {
	tests := []struct {
		name                 string
		key                  []byte
		data                 []byte
		mockUpdateOrderError error
		mockProduceError     error
		expectDLQCall        bool
		expectedDLQTopic     string
	}{
		{
			name:             "Success - Valid JSON and Integration OK",
			key:              []byte("key1"),
			data:             []byte(`{"ordemDeVenda": "123", "etapaAtual": "PENDING"}`),
			mockUpdateOrderError: nil,
			expectDLQCall:    false,
		},
		{
			name:             "Failure - Invalid JSON",
			key:              []byte("key2"),
			data:             []byte(`invalid-json`),
			mockUpdateOrderError: nil,
			expectDLQCall:    true,
			expectedDLQTopic: "dlq-topic",
		},
		{
			name:             "Failure - Integration Error",
			key:              []byte("key3"),
			data:             []byte(`{"ordemDeVenda": "123", "etapaAtual": "PENDING"}`),
			mockUpdateOrderError: errors.New("connection failed"),
			expectDLQCall:    true,
			expectedDLQTopic: "dlq-topic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// flags to verify calls
			dlqCalled := false
			updateOrderCalled := false

			mockProducer := &MockEventProducer{
				ProduceFunc: func(topic string, key []byte, value []byte) error {
					dlqCalled = true
					if topic != tt.expectedDLQTopic {
						t.Errorf("expected DLQ topic %s, got %s", tt.expectedDLQTopic, topic)
					}
					return tt.mockProduceError
				},
			}

			mockIntegration := &MockOrderIntegration{
				UpdateOrderFunc: func(data []byte) error {
					updateOrderCalled = true
					return tt.mockUpdateOrderError
				},
			}

			processor := NewOrderProcessor(mockIntegration, mockProducer, "dlq-topic")

			err := processor.Process(tt.key, tt.data)

			// We expect Process to handle errors gracefully (return nil) 
			// unless it's a critical system failure not covered here.
			if err != nil {
				t.Errorf("Process() error = %v, expected nil", err)
			}

			if tt.expectDLQCall != dlqCalled {
				t.Errorf("expected DLQ called %v, got %v", tt.expectDLQCall, dlqCalled)
			}

			// If JSON was invalid, UpdateOrder should NOT have been called
			if string(tt.data) == "invalid-json" && updateOrderCalled {
				t.Error("expected UpdateOrder NOT to be called for invalid JSON")
			}
		})
	}
}
