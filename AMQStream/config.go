package AMQStream

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Config struct {
	cfg       *kafka.ConfigMap
	consumers []ConsumerOptions
	producers []ProducerOptions
}

type ConsumerOptions struct {
	subscriptions map[string]Subscription
}

type Subscription struct {
	topic       string
	event       ISpecificRecord
	subscriptor ISuscriber
}

type ProducerOptions struct {
	ToPublish map[string][]string
}

type KafkaOption struct {
	BootstrapServers            string
	GroupId                     string
	SessionTimeoutMs            int
	SecurityProtocol            string
	AutoOffsetReset             string
	SslCertificateLocation      string
	MillisecondsTimeout         int
	ConsumerDebug               string
	MaxRetry                    int
	AutoRegisterSchemas         bool
	MessageMaxBytes             int
	PartitionAssignmentStrategy string
	ApplicationName             string
}

type ConsumerMetadata struct {
	Timestamp time.Time
	Key       string
	Header    []kafka.Header
}
