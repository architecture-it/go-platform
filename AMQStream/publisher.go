package AMQStream

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
)

func (c *config) to(event ISpecificRecord, key string) error {

	for _, element := range c.producers {
		for _, topic := range element.ToPublish[event.Schema()] {
			err := c.publish(event, key, topic)
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		}
	}

	return nil
}

func (c *config) publish(event ISpecificRecord, key string, topic string) error {
	appName := getOrDefaultString(configurations, ApplicationName, " ")

	schemaUrl := os.Getenv(SchemaUrl)
	if schemaUrl == "" {
		schemaUrl = configurations[SchemaUrl]
	}

	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaUrl))

	if err != nil {
		log.SugarLogger.Errorf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	serConfig := avro.NewSerializerConfig()

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, serConfig)

	ser.SubjectNameStrategy = WithoutStrategy

	p, err := kafka.NewProducer(c.cfg)

	if err != nil {
		log.SugarLogger.Errorf("Failed to create producer: %s\n", err.Error())
		os.Exit(1)
	}

	deliveryChan := make(chan kafka.Event)

	byteId, _ := json.Marshal(key)

	payload, err := ser.Serialize(event.SchemaName(), event)

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Key:            byteId,
		Timestamp:      time.Time{},
		TimestampType:  0,
		Opaque:         nil,
		Headers:        []kafka.Header{{Key: "remitente", Value: []byte(appName)}},
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Errorf("error: %#v", m.TopicPartition.Error)
		return m.TopicPartition.Error
	}

	fmt.Printf("Delivered %v  to %v", m.Key, m.TopicPartition)
	close(deliveryChan)

	return nil
}

func WithoutStrategy(topic string, serdeType serde.Type, schema schemaregistry.SchemaInfo) (string, error) {
	return topic, nil
}
