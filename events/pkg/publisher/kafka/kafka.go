package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

const (
	clientID           = "creator"
	refreshFrequency   = 30
	createTopicTimeout = 15
)

type Queue struct {
	sink   sarama.SyncProducer
	source sarama.Consumer
}

func New(addresses []string) (*Queue, error) {
	sc := sarama.NewConfig()
	sc.ClientID = clientID

	sc.Version = sarama.V2_4_0_0
	sc.Metadata.RefreshFrequency = refreshFrequency * time.Second
	sc.Producer.Return.Successes = true

	client, err := sarama.NewClient(addresses, sc)
	if err != nil {
		return nil, err
	}

	sink, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	source, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &Queue{
		sink:   sink,
		source: source,
	}, nil
}

func (q *Queue) Publish(payload []byte, topic string) error {
	err := q.ensureKafkaTopic(topic)
	if err != nil {
		return errors.Wrap(err, "creating kafka topic")
	}

	part, offset, err := q.sink.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(payload),
	})
	if err != nil {
		return errors.Wrap(err, "publishing message")
	}

	fmt.Printf("published to partition: %d, position in queue: %d\n", part, offset)

	return nil
}

type Processor func(msg []byte) error

func (q *Queue) Consume(processor Processor, topic string) error {
	consumer, err := q.source.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}

	for m := range consumer.Messages() {
		if err := processor(m.Value); err != nil {
			return err
		}
	}

	return nil
}

func (q *Queue) ensureKafkaTopic(topic string) error {
	broker := sarama.NewBroker("localhost:9092")

	config := sarama.NewConfig()
	config.Version = sarama.V2_4_0_0

	err := broker.Open(config)
	if err != nil {
		return errors.Wrap(err, "opening broker")
	}

	defer broker.Close()

	connected, err := broker.Connected()
	if err != nil {
		return err
	}

	if !connected {
		return errors.New("broked doesnt seem to be connected")
	}

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(1)
	topicDetail.ReplicationFactor = int16(1)
	topicDetail.ConfigEntries = make(map[string]*string)

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[topic] = topicDetail

	request := sarama.CreateTopicsRequest{
		Timeout:      time.Second * createTopicTimeout,
		TopicDetails: topicDetails,
	}

	_, err = broker.CreateTopics(&request)
	if err != nil {
		return errors.Wrap(err, "failed creating topics request")
	}

	return nil
}
