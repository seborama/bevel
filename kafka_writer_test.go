package bevel_test

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/seborama/bevel"
	"github.com/stretchr/testify/require"
)

func TestNewKafkaWriter(t *testing.T) {
	require := require.New(t)

	mockSeedBroker := sarama.NewMockBroker(t, 1)
	defer mockSeedBroker.Close()

	mockSeedBroker.Returns(new(sarama.MetadataResponse))

	kafkaWriter, err := bevel.NewKafkaBEWriter([]string{mockSeedBroker.Addr()}, nil, "test_report_coupon_loaded")
	defer kafkaWriter.Close()
	require.Nil(err)
}

func TestFailedNewKafkaWriter(t *testing.T) {
	require := require.New(t)

	_, err := bevel.NewKafkaBEWriter([]string{"not_a_valid_broken:918723"}, nil, "test_report_coupon_loaded")
	require.NotNil(err)
}

func TestWrite(t *testing.T) {
	const testTopicName = "test_bevel"

	msg := CounterMsg{
		StandardMessage: bevel.StandardMessage{
			EventName:         "test_event",
			CreatedTSUnixNano: time.Now().UnixNano(),
		},
		Counter: 12345,
	}

	require := require.New(t)

	mockSeedBroker := sarama.NewMockBroker(t, 1)
	defer mockSeedBroker.Close()

	s := []struct {
		kafkaErr sarama.KError
	}{
		{ // Happy path
			kafkaErr: sarama.ErrNoError,
		},
		{ // Sad path
			kafkaErr: sarama.ErrInvalidMessage,
		},
	}

	for _, td := range s {
		metadataResponse := new(sarama.MetadataResponse)
		metadataResponse.AddBroker(mockSeedBroker.Addr(), mockSeedBroker.BrokerID())
		metadataResponse.AddTopicPartition(testTopicName, 0, mockSeedBroker.BrokerID(), nil, nil, sarama.ErrNoError)
		mockSeedBroker.Returns(metadataResponse) // connect producer

		prodResponse := new(sarama.ProduceResponse)
		prodResponse.AddTopicPartition(testTopicName, 0, td.kafkaErr)
		mockSeedBroker.Returns(prodResponse) // publish msg

		config := sarama.NewConfig()
		config.Metadata.Retry.Max = 0
		config.Producer.Retry.Max = 0

		kafkaWriter, err := bevel.NewKafkaBEWriter([]string{mockSeedBroker.Addr()}, config, testTopicName)
		defer kafkaWriter.Close()

		require.Nil(err)

		err = kafkaWriter.Write(msg)
		if td.kafkaErr == sarama.ErrNoError {
			require.Nil(err)
		} else {
			require.Error(err)
		}
	}
}
