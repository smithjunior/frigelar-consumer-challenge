#!/bin/bash

# Configuration
CONTAINER_NAME="kafka"
BOOTSTRAP_SERVER="localhost:29092"
TOPIC="orders"

ORDER_ID="order12345"
STATUS="FATURADO"

PAYLOAD="{\"ordemDeVenda\": \"$ORDER_ID\", \"etapaAtual\": \"$STATUS\"}"

echo "Sending message to topic '$TOPIC'..."
echo "Payload: $PAYLOAD"


echo "$PAYLOAD" | docker exec -i $CONTAINER_NAME kafka-console-producer \
    --bootstrap-server $BOOTSTRAP_SERVER \
    --topic $TOPIC

if [ $? -eq 0 ]; then
    echo "✅ Message sent successfully!"
else
    echo "❌ Failed to send message."
fi
