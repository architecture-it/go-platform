package AMQStream

import "github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"

func To(event avro.SpecificAvroMessage, key string) error {

	return getInstance().to(event, key)
}
