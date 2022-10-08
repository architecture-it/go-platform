package AMQStream

import (
	"fmt"
	"time"

	logger "github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/exp/slices"
)

func (c *config) to(event ISpecificRecord, key, exactopic string, headers []kafka.Header) error {

	// Select the topic
	if exactopic != "" {
		err := c.publish(event, key, exactopic, headers)
		if err != nil {
			logger.SugarLogger.Errorln(err.Error())
			return err
		}
		return nil
	}

	for _, element := range c.producers {
		for _, topic := range element.ToPublish[event.Schema()] {
			err := c.publish(event, key, topic, headers)
			if err != nil {
				logger.SugarLogger.Errorln(err.Error())
				return err
			}
		}
	}

	return nil
}

func (c *config) publish(event ISpecificRecord, key string, topic string, headers []kafka.Header) error {
	// set default Header

	currentRetry := slices.IndexFunc(headers, func(x kafka.Header) bool {
		return x.Key == Remitente
	})

	if currentRetry == -1 {
		appName := getOrDefaultString(configurations, ApplicationName, " ")
		headers = append(headers, kafka.Header{Key: Remitente, Value: []byte(appName)})
	}

	p, err := kafka.NewProducer(c.cfgProducer)

	if err != nil {
		logger.SugarLogger.Errorln(fmt.Sprintf("Failed to create producer: %s\n", err))
		return err
	}
	value, err := serializeMessage(c, event)
	if err != nil || len(value) == 0 {
		logger.SugarLogger.Errorln(fmt.Sprintf("Serialize error: %s\n", err))
		return err
	}

	return producerMessage(c, p, value, key, topic, headers)

}

func producerMessage(c *config,
	p *kafka.Producer,
	value []byte,
	key string,
	topic string,
	headers []kafka.Header) error {

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				logger.SugarLogger.Errorln(fmt.Sprintf("Error: %v\n", ev))
			default:
				logger.SugarLogger.Warnln(fmt.Sprintf("Ignored event: %s\n", ev))
			}
		}
	}()
	// A delivery channel for each message sent.
	// This permits to receive delivery reports
	// separately and to handle the use case
	// of a server that has multiple concurrent
	// produce requests and needs to deliver the replies
	// to many different response channels.
	deliveryChan := make(chan kafka.Event)
	go func() {
		for e := range deliveryChan {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or
				// permanent failure after retries have been exhausted.
				// Application level retries won't help since the client
				// is already configured to do that.
				m := ev
				if m.TopicPartition.Error != nil {
					logger.SugarLogger.Infoln(fmt.Sprintf("Delivery failed: %v\n", m.TopicPartition.Error))
				} else {
					logger.SugarLogger.Infoln(fmt.Sprintf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset))
				}

			default:
				logger.SugarLogger.Debugln(fmt.Sprintf("Ignored event: %s\n", ev))
			}
			// in this case the caller knows that this channel is used only
			// for one Produce call, so it can close it.
			close(deliveryChan)
		}
	}()

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
		Headers:        headers,
		Key:            []byte(key),
		Timestamp:      time.Time{},
		TimestampType:  0,
		Opaque:         nil,
	}, deliveryChan)

	if err != nil {
		close(deliveryChan)
		if err.(kafka.Error).Code() == kafka.ErrQueueFull {
			// Producer queue is full, wait 1s for messages
			// to be delivered then try again.
			time.Sleep(time.Second)
		}
		logger.SugarLogger.Errorln(fmt.Sprintf("Failed to produce message: %v\n", err))
	}

	// Flush and close the producer and the events channel
	for p.Flush(10000) > 0 {
		logger.SugarLogger.Debugf("Still waiting to flush outstanding messages\n", err)
	}
	p.Close()
	return nil
}
