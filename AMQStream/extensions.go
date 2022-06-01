package AMQStream

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/architecture-it/ARQ.Common-SettingsGO/extensions"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mitchellh/mapstructure"
)

var configurations = make(map[string]string)

func AddKafka() (*Config, error) {
	config, err := bindConfiguration()
	if err != nil {
		panic(err)
	}

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers":                   config.BootstrapServers,
		"group.id":                            config.GroupId,
		"security.protocol":                   config.SecurityProtocol,
		"ssl.certificate.location":            config.SslCertificateLocation,
		"message.max.bytes":                   config.MessageMaxBytes,
		"enable.ssl.certificate.verification": false,
		// "application.id": config.ApplicationName,
	}
	return &Config{cfg: kafkaConfig}, nil
}

func bindConfiguration() (*KafkaOption, error) {
	configuration := extensions.GetConfiguration("enviroment.yaml")
	mapstructure.Decode(configuration["Kafka"], &configurations)

	messageMaxBytes, errorConversion := strconv.Atoi(configurations[MessageMaxBytes])

	if errorConversion != nil {

		return nil, errorConversion
	}

	err := validRequired()

	if err != nil {

		return nil, err
	}

	result := KafkaOption{
		BootstrapServers:            getOrDefaultString(BootstrapServers, configurations[BootstrapServers]),
		GroupId:                     getOrDefaultString(GroupId, configurations[GroupId]),
		SessionTimeoutMs:            getOrDefaultInt(SessionTimeoutMs, 60000),
		SecurityProtocol:            getOrDefaultString(SecurityProtocol, configurations[SecurityProtocol]),
		AutoOffsetReset:             getOrDefaultString(AutoOffsetReset, "Earlitest"),
		SslCertificateLocation:      getOrDefaultString(SslCertificateLocation, configurations[SslCertificateLocation]),
		MillisecondsTimeout:         getOrDefaultInt(MillisecondsTimeout, 10000),
		ConsumerDebug:               getOrDefaultString(ConsumerDebug, ""),
		MaxRetry:                    getOrDefaultInt(MaxRetry, 3),
		AutoRegisterSchemas:         getOrDefaultBool(AutoRegisterSchemas, true),
		MessageMaxBytes:             getOrDefaultInt(MessageMaxBytes, messageMaxBytes),
		PartitionAssignmentStrategy: getOrDefaultString(ConsumerDebug, "CooperativeSticky"),
		// ApplicationName: getOrDefaultString(ApplicationName, configurations[ApplicationName]),
	}
	return &result, nil
}

func validRequired() error {
	boopstrapServer := os.Getenv(BootstrapServers)
	if boopstrapServer == "" && configurations[BootstrapServers] == "" {
		return errors.New("The BootstrapServer is requiered")
	}
	applicationName := os.Getenv(ApplicationName)
	if applicationName == "" && configurations[ApplicationName] == "" {
		return errors.New("The ApplicationName is requiered")
	}
	schemaUrl := os.Getenv(SchemaUrl)
	if schemaUrl == "" && configurations[SchemaUrl] == "" {
		return errors.New("The SchemaUrl is requiered")
	}
	return nil
}

func (c *Config) ToConsumer(suscriber ISuscriber, event ISpecificRecord, topic string) {
	subscriptions := make(map[string]Subscription)

	subcription := Subscription{
		event:       event,
		topic:       topic,
		subscriptor: suscriber,
	}
	subscriptions[event.SchemaName()] = subcription

	c.consumers = append(c.consumers, ConsumerOptions{
		subscriptions: subscriptions,
	})
}

func (c *Config) ToProducer(event ISpecificRecord, topics []string) {
	appended := false

	for _, v := range c.producers {
		if v.ToPublish[event.SchemaName()] != nil {
			for _, t := range topics {
				v.ToPublish[event.SchemaName()] = append(v.ToPublish[event.SchemaName()], t)
			}
			appended = true
		}
	}
	if !appended {
		topicsForAdd := make(map[string][]string)
		topicsForAdd[event.SchemaName()] = topics
		c.producers = append(c.producers, ProducerOptions{
			ToPublish: topicsForAdd,
		})
	}
}

func (c *Config) Build() {
	for index, element := range c.consumers {
		for indexj, suscriber := range element.subscriptions {
			go func() error {
				for {
					err := c.Consumer(suscriber.event, suscriber.topic)
					if err != nil {
						return err
					}
				}

				return nil
			}()
			fmt.Println(indexj)
		}

		fmt.Println(index)
	}
}
