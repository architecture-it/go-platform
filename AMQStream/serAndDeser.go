package AMQStream

import (
	"fmt"

	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
)

func serializeMessage(c *config, event ISpecificRecord) ([]byte, error) {
	client, err := schemaregistry.NewClient(c.schemaRegistry)

	if err != nil {
		log.SugarLogger.Errorln("Failed to create schema registry client: %s\n", err)
		return nil, err
	}

	serConfig := avro.NewSerializerConfig()

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, serConfig)

	if err != nil {
		log.SugarLogger.Errorln("Failed to serializer: %s\n", err)
		return nil, err
	}
	ser.SubjectNameStrategy = withoutStrategy

	return ser.Serialize(event.SchemaName(), event)

}

func withoutStrategy(topic string, serdeType serde.Type, schema schemaregistry.SchemaInfo) (string, error) {
	return topic, nil
}

func deserializeMessage(c *config, message *kafka.Message, event ISpecificRecord) (ISpecificRecord, error) {
	client, err := schemaregistry.NewClient(c.schemaRegistry)

	if event == nil {
		errorMessage := fmt.Sprintf("Event is nil. Key: %s\n", message.Key)
		log.SugarLogger.Errorln(errorMessage)
		return nil, fmt.Errorf(errorMessage)
	}

	if err != nil {
		log.SugarLogger.Errorln(fmt.Sprintf("Failed to create schema registry client: %s\n", err))
		return nil, err
	}

	deser, err := avro.NewSpecificDeserializer(client, serde.ValueSerde, avro.NewDeserializerConfig())

	if err != nil {
		log.SugarLogger.Errorln("Failed to deserializer: %s\n", err)
		return nil, err
	}
	deser.SubjectNameStrategy = withoutStrategy

	result := event

	err = deser.DeserializeInto(event.SchemaName(), message.Value, result)
	if err != nil {
		log.SugarLogger.Errorln(fmt.Sprintf("Failed to deserialize payload: %s\n", err))
		return nil, err
	}

	return result, nil
}
