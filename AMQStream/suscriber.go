package AMQStream

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/mitchellh/mapstructure"
)

func (c *Config) NotifyToSubscriber(event ISpecificRecord, topic string, metadata ConsumerMetadata) error {

	for index, element := range c.consumers {
		subscription := element.subscriptions[event.SchemaName()]

		if subscription.topic == topic {
			subscription.subscriptor.Handler(event, metadata)
		}

		fmt.Println(index)
	}

	return nil
}

func (k *Config) Consumer(event ISpecificRecord, topic string) error {
	eventSchema := event.Schema()

	var result []byte
	var metadata ConsumerMetadata

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println(k.cfg)

	c, err := kafka.NewConsumer(k.cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	var topics []string
	topics = append(topics, topic)

	err = c.SubscribeTopics(topics, nil)

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(30)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
					metadata.Header = e.Headers
					metadata.Key = string(e.Key)
					metadata.Timestamp = e.Timestamp
				}

				result = e.Value
				run = false
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()

	codec, errr := goavro.NewCodec(eventSchema)

	if errr != nil {
		fmt.Println(errr)
	}

	decoded, _, errr := codec.NativeFromBinary(result[5:])
	if errr != nil {
		fmt.Println(errr)
	}

	fmt.Println(fmt.Sprintf("%s", decoded))

	mapstructure.Decode(decoded, &event)

	resultEvent := k.NotifyToSubscriber(event, topic, metadata)

	fmt.Println(resultEvent)

	return nil
}
