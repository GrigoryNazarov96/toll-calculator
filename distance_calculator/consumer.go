package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/GrigoryNazarov96/toll-calculator/aggregator/client"
	"github.com/GrigoryNazarov96/toll-calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService Calculator
	client      client.Client
}

func NewKafkaConsumer(topic string, s Calculator, cl client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}
	c.SubscribeTopics([]string{topic}, nil)
	return &KafkaConsumer{
		consumer:    c,
		calcService: s,
		client:      cl,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka consumer started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) Close() {
	c.isRunning = false
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consumer error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("serialization error %s", err)
			continue
		}
		distance := c.calcService.CalculateDistance(data)
		telemetryData := &types.TelemetryDataRequest{
			Distance: distance,
			ObuID:    int32(data.OBUID),
			Unix:     time.Now().Unix(),
		}
		if err := c.client.Aggregate(context.Background(), telemetryData); err != nil {
			logrus.Errorf("aggregate error: %s", err)
			continue
		}
	}
}
