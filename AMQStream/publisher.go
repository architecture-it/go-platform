package AMQStream

import (
	"encoding/json"
	"os"
	"time"

	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
)

func (c *config) to(event ISpecificRecord, key string) error {

	for _, element := range c.producers {
		for _, topic := range element.ToPublish[event.SchemaName()] {
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
	eventBytes, _ := event.MarshalJSON()
	eventSchema := event.Schema()

	p, err := kafka.NewProducer(c.cfg)

	if err != nil {
		log.SugarLogger.Errorf("Failed to create producer: %s\n", err.Error())
		os.Exit(1)
	}

	deliveryChan := make(chan kafka.Event)

	byteId, _ := json.Marshal(key)

	codec, err := goavro.NewCodec(eventSchema)

	if err != nil {
		return err
	}

	native, _, err := codec.NativeFromTextual(eventBytes)
	if err != nil {
		return err
	}

	var bin []byte
	bin = append(bin, 0)
	bin = append(bin, 0)
	bin = append(bin, 0)
	bin = append(bin, 0)
	bin = append(bin, 223)
	binary, err := codec.BinaryFromNative(nil, native)
	if err != nil {
		return err
	}

	for _, element := range binary {
		bin = append(bin, element)
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bin,
		Key:            byteId,
		Timestamp:      time.Time{},
		TimestampType:  0,
		Opaque:         nil,
		Headers:        []kafka.Header{{Key: "remitente", Value: []byte("AMQTest")}},
	}, deliveryChan)

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	close(deliveryChan)

	return nil
}
