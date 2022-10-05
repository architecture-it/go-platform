package AMQStream

import (
	"errors"
	"os"

	extension "github.com/architecture-it/go-platform/config"
	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mitchellh/mapstructure"
)

var configurations = make(map[string]string)

func AddKafka() (*config, error) {
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
		"enable.ssl.certificate.verification": config.EnableSslCertificateVerification,
		"auto.offset.reset":                   config.AutoOffsetReset,
		"session.timeout.ms":                  config.SessionTimeoutMs,
		"debug":                               config.ConsumerDebug,
		"partition.assignment.strategy":       config.PartitionAssignmentStrategy,
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
		BootstrapServers:                 getOrDefaultString(configurations, BootstrapServers, configurations[BootstrapServers]),
		GroupId:                          getOrDefaultString(configurations, GroupId, ""),
		SessionTimeoutMs:                 getOrDefaultInt(configurations, SessionTimeoutMs, 60000),
		SecurityProtocol:                 getOrDefaultString(configurations, SecurityProtocol, "plaintext"),
		AutoOffsetReset:                  getOrDefaultString(configurations, AutoOffsetReset, "earliest"),
		SslCertificateLocation:           getOrDefaultString(configurations, SslCertificateLocation, ""),
		MillisecondsTimeout:              getOrDefaultInt(configurations, MillisecondsTimeout, 10000),
		ConsumerDebug:                    getOrDefaultString(configurations, ConsumerDebug, " "),
		MaxRetry:                         getOrDefaultInt(configurations, MaxRetry, 3),
		PartitionAssignmentStrategy:      getOrDefaultString(configurations, PartitionAssignmentStrategy, "cooperative-sticky"),
		MessageMaxBytes:                  getOrDefaultInt(configurations, MessageMaxBytes, 100000),
		EnableSslCertificateVerification: getOrDefaultBool(configurations, EnableSslCertificateVerification, false),
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

func (c *config) ToConsumer(suscriber ISuscriber, event ISpecificRecord, topic string) {
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

func (c *config) ToProducer(event ISpecificRecord, topics []string) {
	appended := false

	for _, v := range c.producers {
		if v.ToPublish[event.Schema()] != nil {
			for _, t := range topics {
				v.ToPublish[event.Schema()] = append(v.ToPublish[event.Schema()], t)
			}
			appended = true
		}
	}
	if !appended {
		topicsForAdd := make(map[string][]string)
		topicsForAdd[event.Schema()] = topics
		c.producers = append(c.producers, ProducerOptions{
			ToPublish: topicsForAdd,
		})
	}
}

func (c *config) Build() {
	for _, element := range c.consumers {
		for _, suscriber := range element.subscriptions {
			go func() error {
				for {
					event := suscriber.event
					topic := suscriber.topic
					err := c.consumer(event, topic)
					if err != nil {
						log.Logger.Error(err.Error())
						return err
					}
				}

			}()
		}

	}

}
