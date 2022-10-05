package AMQStream

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/linkedin/goavro/v2"
	"github.com/mitchellh/mapstructure"
)

func (c *config) NotifyToSubscriber(event interface{}, topic string, metadata ConsumerMetadata) error {

	for _, element := range c.consumers {
		subscription := element.subscriptions[event.SchemaName()]

		if subscription.topic == topic {
			subscription.subscriptor.Handler(event, metadata)
		}

	}

	return nil
}

func (k *config) consumer(event interface{}, topic string) error {
	eventSchema := event.Schema()
	millisecondTimeout := getOrDefaultInt(configurations, MillisecondsTimeout, 10000)

	var result []byte
	var metadata ConsumerMetadata

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	log.SugarLogger.Infof("%v", k.cfg)

	c, err := kafka.NewConsumer(k.cfg)

	if err != nil {
		log.SugarLogger.Errorf("Failed to create consumer: %s\n", err.Error())
		os.Exit(1)
	}

	var topics []string
	topics = append(topics, topic)

	err = c.SubscribeTopics(topics, nil)

	if err != nil {
		return nil
	}
	run := true

	for run {
		select {
		case sig := <-sigchan:
			log.SugarLogger.Infof("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := c.Poll(millisecondTimeout)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				if e.Headers != nil {
					metadata.Header = e.Headers
					metadata.Key = string(e.Key)
					metadata.Timestamp = e.Timestamp
				}

				result = e.Value
				codec, errr := goavro.NewCodec(eventSchema)

				if errr != nil {
					return errr
				}

				decoded, _, errr := codec.NativeFromBinary(result[5:])
				if errr != nil {
					return errr
				}

				mapstructure.Decode(decoded, &event)

				k.NotifyToSubscriber(event, topic, metadata)
			case kafka.Error:
				log.SugarLogger.Errorf("%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.SugarLogger.Infof("Ignored %v\n", e)
			}
		}
	}

	c.Close()

	return nil
}
