package AMQStream

import "github.com/confluentinc/confluent-kafka-go/kafka"

func To(event ISpecificRecord, key string) error {

	return getInstance().to(event, key, "", nil)
}

func ToWithHeader(event ISpecificRecord, key string, headers []kafka.Header) error {

	return getInstance().to(event, key, "", headers)
}

func ToInTopic(event ISpecificRecord, key, exactopic string) error {

	return getInstance().to(event, key, exactopic, nil)
}

func ToInTopicWithHeader(event ISpecificRecord, key, exactopic string, headers []kafka.Header) error {

	return getInstance().to(event, key, exactopic, headers)
}
