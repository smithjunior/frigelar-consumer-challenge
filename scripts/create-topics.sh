#!/bin/bash

CONTAINER_NAME="kafka"
PARTITIONS=1
REPLICATION_FACTOR=1
BOOTSTRAP_SERVER="localhost:29092"

TOPICS=("orders" "orders-dlq")

echo "Waiting for Kafka to be ready..."
sleep 5

for TOPIC in "${TOPICS[@]}"; do
    echo "Creating topic: $TOPIC"
    docker exec $CONTAINER_NAME kafka-topics --create \
        --if-not-exists \
        --bootstrap-server $BOOTSTRAP_SERVER \
        --replication-factor $REPLICATION_FACTOR \
        --partitions $PARTITIONS \
        --topic $TOPIC
    
    if [ $? -eq 0 ]; then
        echo "Successfully checked/created topic: $TOPIC"
    else
        echo "Failed to create topic: $TOPIC"
    fi
done

echo "Listing current topics:"
docker exec $CONTAINER_NAME kafka-topics --list --bootstrap-server $BOOTSTRAP_SERVER
