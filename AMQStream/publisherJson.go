package AMQStream

import (
	"encoding/json"
	"fmt"

	logger "github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (c *config) toJson(event ISpecificRecord, metadata ConsumerMetadata) error {
	var headers []kafka.Header

	p, err := kafka.NewProducer(c.cfgProducer)
	if err != nil {
		logger.SugarLogger.Errorln(fmt.Sprintf("Failed to create producer: %s\n", err))
		return err
	}
	appName := getOrDefaultString(configurations, ApplicationName, " ")
	headers = append(headers, kafka.Header{Key: Remitente, Value: []byte(appName)})

	key := metadata.Key
	topic := CrossDeadline
	valueSerialize, _ := serializeMessage(c, topic, event)
	dict := make(map[string]string)

	for i := range metadata.Header {
		dict[metadata.Header[i].Key] = string(metadata.Header[i].Value)
	}
	value := DeadlineMessage{
		Payload:          event,
		PayloadSerialize: valueSerialize,
		ApplicationName:  appName,
		Key:              metadata.Key,
		Timestamp:        metadata.Timestamp,
		Headers:          dict,
		Partition:        int(metadata.Partition),
		Offset:           metadata.Offset,
		Topic:            metadata.Topic,
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		logger.SugarLogger.Errorln(fmt.Sprintf("Can't serialize: %s\n", err))
		return err
	}
	return producerMessage(c, p, bytes, key, topic, headers)
}
