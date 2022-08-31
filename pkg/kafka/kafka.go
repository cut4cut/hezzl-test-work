package producer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
)

const (
	_defaultBootstrapServers = "kafka:29092"
)

type KafkaProducer struct {
	sarama.SyncProducer
}

func New() (*KafkaProducer, error) {
	brokerList := _defaultBootstrapServers

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(strings.Split(brokerList, ","), config)
	if err != nil {
		return nil, fmt.Errorf("failed to open kafka producer, error: %w", err)
	}

	return &KafkaProducer{p}, nil
}

func (p *KafkaProducer) produce(messages string, topicName string) error {
	message := &sarama.ProducerMessage{Topic: topicName, Partition: int32(-1)}
	message.Value = sarama.StringEncoder(string(messages))

	_, _, err := p.SendMessage(message)
	if err != nil {
		return fmt.Errorf("kafka producer can't connect, error: %w", err)
	}

	return nil
}

func (p *KafkaProducer) ProduceStruct(str interface{}, topicName string) error {
	stringJson, err := json.Marshal(str)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	p.produce(string(stringJson), topicName)

	return nil
}
