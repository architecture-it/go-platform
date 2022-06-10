package AMQStream

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
)

func (c *Config) To(event ISpecificRecord, key string) error {

	for index, element := range c.producers {
		for indej, topic := range element.ToPublish[event.SchemaName()] {
			err := c.Publish(event, key, topic)
			if err != nil {
				return err
			}
			fmt.Println(indej)
		}
		fmt.Println(index)
	}

	return nil
}

func (c *Config) Publish(event ISpecificRecord, key string, topic string) error {
	eventBytes, err := event.MarshalJSON()
	eventSchema := event.Schema()

	p, err := kafka.NewProducer(c.cfg)

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	deliveryChan := make(chan kafka.Event)

	//Generate Guid identity
	// id := uuid.New()
	// uuid := strings.Replace(id.String(), "-", "", -1)
	byteId, err := json.Marshal(key)
	fmt.Println(byteId)

	codec, err := goavro.NewCodec(eventSchema)

	if err != nil {
		return err
		// fmt.Println(errr)
	}

	fmt.Println(codec)

	native, _, err := codec.NativeFromTextual(eventBytes)
	if err != nil {
		// fmt.Println(err)
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

	for index, element := range binary {
		bin = append(bin, element)
		fmt.Println(index)
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
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)

	return nil
}
