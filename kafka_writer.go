package bevel

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/Shopify/sarama"
)

// ErrUnableToSendMessage occurs when the message could not be
// sent to the external queue (Apache Kafka, etc).
var ErrUnableToSendMessage = errors.New("unable to send message")

// KafkaBEWriter is an Apache Kafka mock Writer implementation..
type KafkaBEWriter struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaBEWriter creates a new KafkaBEWriter using the given broker addresses and configuration.
func NewKafkaBEWriter(addrs []string, config *sarama.Config, topic string) Writer {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return &KafkaBEWriter{
		producer: producer,
		topic:    topic,
	}
}

// Write outputs the contents of Message to a Kafka topic.
func (bew *KafkaBEWriter) Write(m Message) error {
	// TODO - change the hard-coding of json.Marshal() to a
	// strategy pattern.
	// This would allow to inject a strategy to write messages.
	// For instance, a JSONWriterStrategy or a MapWriterStrategy (which
	// would return a map) or an XMLWriterStrategy, etc.
	jsonStr, err := json.Marshal(m)
	if err != nil {
		return ErrUnableToTranscodeMessage
	}

	msg := &sarama.ProducerMessage{Topic: bew.topic, Value: sarama.StringEncoder(jsonStr)}
	_, _, err = bew.producer.SendMessage(msg)
	if err != nil {
		return ErrUnableToSendMessage
	}

	return nil
}

// Close shuts down the producer and flushes any messages it may have buffered.
// You must call this function before a producer object passes out of scope, as
// it may otherwise leak memory. You must call this before calling Close on the
// underlying client.
func (bew *KafkaBEWriter) Close() error {
	err := bew.producer.Close()
	return err
}
