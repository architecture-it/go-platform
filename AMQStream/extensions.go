package AMQStream

import (
	"errors"
	"fmt"
	"os"

	extension "github.com/architecture-it/go-platform/config"
	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mitchellh/mapstructure"
)

var configurations = make(map[string]string)

func AddKafka() (*Config, error) {
	config, err := bindConfiguration()
	if err != nil {
		log.Logger.DPanic(err.Error())
		panic(err)
	}

	getInstance().cfg = &kafka.ConfigMap{
		"bootstrap.servers":                   config.BootstrapServers,
		"group.id":                            config.GroupId,
		"security.protocol":                   config.SecurityProtocol,
		"ssl.certificate.location":            config.SslCertificateLocation,
		"message.max.bytes":                   config.MessageMaxBytes,
		"enable.ssl.certificate.verification": false,
		"auto.offset.reset":                   config.AutoOffsetReset,
	}

	return getInstance(), nil
}

func bindConfiguration() (*KafkaOption, error) {
	configuration := extension.GetConfiguration("enviroment.yaml")
	mapstructure.Decode(configuration["Kafka"], &configurations)

	err := validRequired()

	if err != nil {

		return nil, err
	}

	result := KafkaOption{
		BootstrapServers:            getOrDefaultString(configurations, BootstrapServers, configurations[BootstrapServers]),
		GroupId:                     getOrDefaultString(configurations, GroupId, ""),
		SessionTimeoutMs:            getOrDefaultInt(configurations, SessionTimeoutMs, 60000),
		SecurityProtocol:            getOrDefaultString(configurations, SecurityProtocol, "plaintext"),
		AutoOffsetReset:             getOrDefaultString(configurations, AutoOffsetReset, "earliest"),
		SslCertificateLocation:      getOrDefaultString(configurations, SslCertificateLocation, ""),
		MillisecondsTimeout:         getOrDefaultInt(configurations, MillisecondsTimeout, 10000),
		ConsumerDebug:               getOrDefaultString(configurations, ConsumerDebug, ""),
		MaxRetry:                    getOrDefaultInt(configurations, MaxRetry, 3),
		AutoRegisterSchemas:         getOrDefaultBool(configurations, AutoRegisterSchemas, true),
		PartitionAssignmentStrategy: getOrDefaultString(configurations, ConsumerDebug, "CooperativeSticky"),
		MessageMaxBytes:             getOrDefaultInt(configurations, MessageMaxBytes, 100000),
		// ApplicationName: getOrDefaultString(ApplicationName, configurations[ApplicationName]),
	}
	return &result, nil
}

func validRequired() error {
	boopstrapServer := os.Getenv(BootstrapServers)
	if boopstrapServer == "" && configurations[BootstrapServers] == "" {
		return errors.New("the bootstrapServer is requiered")
	}
	applicationName := os.Getenv(ApplicationName)
	if applicationName == "" && configurations[ApplicationName] == "" {
		return errors.New("the applicationName is requiered")
	}
	schemaUrl := os.Getenv(SchemaUrl)
	if schemaUrl == "" && configurations[SchemaUrl] == "" {
		return errors.New("the schemaUrl is requiered")
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
						log.Logger.Error(err.Error())
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
