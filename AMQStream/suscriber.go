package AMQStream

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/exp/slices"
)

func (c *config) notifyToSubscriber(event ISpecificRecord, metadata ConsumerMetadata) error {
	for _, subscription := range c.consumers {
		if slices.Contains(subscription.topic, metadata.Topic) {
			err := subscription.subscriptor.Handler(event, metadata)
			if err != nil {
				moveRetryTopic(&subscription, event, metadata, err)
			}
			return nil
		}
	}
	return errors.New("the message could not be notified")
}

func (k *config) consumer(event ISpecificRecord, topic []string, wg *sync.WaitGroup) error {
	defer wg.Done()
	millisecondTimeout := getOrDefaultInt(configurations, MillisecondsTimeout, 10000)
	lastOffset := []kafka.TopicPartition{}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(k.cfgConsumer)

	if err != nil {
		log.SugarLogger.Errorf("Failed to create consumer: %s\n", err.Error())
		os.Exit(1)
	}
	err = c.SubscribeTopics(topic, nil)
	// deserialize
	deser, err := createDeserialize(k)
	if err != nil {
		return nil
	}
	run := true
	// defer closeConsumer(c, topic)
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
				if compareTopicPartitionOffset(e, lastOffset) {
					log.SugarLogger.Warnln(fmt.Sprintf("Consume offset duplicated | Consumed %v", e.TopicPartition.Offset))
					continue
				}
				createOrUpdateStore(e, &lastOffset)

				log.SugarLogger.Infof(fmt.Sprintf("Topic: %v | Consume offset: %v| partition: %v", *e.TopicPartition.Topic, e.TopicPartition.Offset, e.TopicPartition.Partition))
				metadata := createMetadata(e)
				eventDes, err := deserializeMessage(deser, e, event)
				if err != nil {
					continue
				}
				if filterMessage(e, k.MaxRetry) {
					_ = publishDeadline(k, eventDes, metadata)
					continue
				}

				k.notifyToSubscriber(eventDes, metadata)

			case kafka.Error:
				log.SugarLogger.Errorf(fmt.Sprintf("%% Error: %v: %v\n", e.Code(), e))
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				log.SugarLogger.Debugln(fmt.Sprintf("Ignored %v\n", e))
			}
		}
	}
	closeConsumer(c, topic)
	return nil
}

func createMetadata(msg *kafka.Message) ConsumerMetadata {
	var metadata ConsumerMetadata
	if msg.Headers != nil {
		metadata.Header = msg.Headers
	}
	metadata.Key = string(msg.Key)
	metadata.Timestamp = msg.Timestamp
	metadata.Topic = *msg.TopicPartition.Topic
	metadata.Offset = int64(msg.TopicPartition.Offset)
	metadata.Partition = msg.TopicPartition.Partition
	return metadata
}

func filterMessage(msg *kafka.Message, maxRetry int) bool {
	currentRetry := -1
	for i, x := range msg.Headers { // the last index
		if x.Key == RetryCount {
			currentRetry = i
		}
	}
	deadline := slices.IndexFunc(msg.Headers, func(x kafka.Header) bool {
		return x.Key == Deadline
	})
	if currentRetry != -1 {
		retry, _ := strconv.Atoi(string(msg.Headers[currentRetry].Value))
		return retry >= maxRetry ||
			deadline != -1
	}

	return false
}

func publishDeadline(k *config, event ISpecificRecord, metadata ConsumerMetadata) error {
	deadline := slices.IndexFunc(metadata.Header, func(x kafka.Header) bool {
		return x.Key == Deadline
	})

	if deadline == -1 {
		metadata.Header = append(metadata.Header, kafka.Header{Key: Deadline, Value: []byte("true")})
		return k.toJson(event, metadata)
	}

	return nil
}

func moveRetryTopic(configuration *ConsumerOptions, event ISpecificRecord, metadata ConsumerMetadata, err error) {
	if configuration.move != "" {
		log.SugarLogger.Errorln(fmt.Sprintf("an error has occurred. The message is moved to the topic %v", configuration.move))

		retry := "1"
		currentRetry := -1
		for i, x := range metadata.Header { // the last index
			if x.Key == RetryCount {
				currentRetry = i
			}
		}

		if currentRetry != -1 {
			aux, _ := strconv.Atoi(string(metadata.Header[currentRetry].Value))
			retry = strconv.Itoa(aux + 1)
		}

		metadata.Header = append(metadata.Header, kafka.Header{Key: RetryCount, Value: []byte(retry)})
		metadata.Header = append(metadata.Header, kafka.Header{Key: fmt.Sprintf(Reason, retry), Value: []byte(err.Error())})

		ToInTopicWithHeader(event, metadata.Key, configuration.move, metadata.Header)
	}
}

func closeConsumer(c *kafka.Consumer, topic []string) {
	log.SugarLogger.Infoln(fmt.Sprintf("Close consumer %v", topic))
	c.Close()
}

func compareTopicPartitionOffset(message *kafka.Message, lastOffset []kafka.TopicPartition) bool {
	compare := message.TopicPartition
	for _, v := range lastOffset {
		if *v.Topic == *compare.Topic && v.Partition == compare.Partition {
			return v.Offset >= compare.Offset
		}
	}
	return false
}

func createOrUpdateStore(message *kafka.Message, lastOffset *[]kafka.TopicPartition) {
	compare := message.TopicPartition
	for i, v := range *lastOffset {
		if *v.Topic == *compare.Topic && v.Partition == compare.Partition {
			(*lastOffset)[i] = message.TopicPartition
			return
		}
	}
	(*lastOffset) = append((*lastOffset), message.TopicPartition)
}
