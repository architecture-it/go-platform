package AMQStream

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
)

type config struct {
	cfgProducer    *kafka.ConfigMap
	cfgConsumer    *kafka.ConfigMap
	schemaRegistry *schemaregistry.Config
	consumers      []ConsumerOptions
	producers      []ProducerOptions
	MaxRetry       int
}

type ConsumerOptions struct {
	topic       []string
	event       ISpecificRecord
	subscriptor ISuscriber
	move        string
}

type Subscription struct {
}

type ProducerOptions struct {
	ToPublish map[string][]string
}

type KafkaOption struct {
	BootstrapServers                 string
	SchemaRegistry                   string
	GroupId                          string
	SessionTimeoutMs                 int
	SecurityProtocol                 string
	AutoOffsetReset                  string
	SslCertificateLocation           string
	MillisecondsTimeout              int
	ConsumerDebug                    string
	MaxRetry                         int
	AutoRegisterSchemas              bool
	MessageMaxBytes                  int
	PartitionAssignmentStrategy      string
	ApplicationName                  string
	EnableSslCertificateVerification bool
}

type ConsumerMetadata struct {
	Timestamp time.Time
	Key       string
	Header    []kafka.Header
	Topic     string
	Partition int32
	Offset    int64
}

type DeadlineMessage struct {
	Payload          interface{}
	PayloadSerialize []byte
	ApplicationName  string
	Key              string
	Timestamp        time.Time
	Headers          map[string]string
	Partition        int
	Offset           int64
	Topic            string
}
