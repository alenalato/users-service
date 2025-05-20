#!/bin/bash

/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic "$TOPIC_NAME" --bootstrap-server kafka:9092
echo "topic $TOPIC_NAME created"
