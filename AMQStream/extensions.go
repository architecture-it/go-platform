package AMQStream

import (
	"errors"
	"os"
	"sync"

	extension "github.com/architecture-it/go-platform/config"
	"github.com/architecture-it/go-platform/log"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/mitchellh/mapstructure"
)

var configurations = make(map[string]string)

func AddKafka() (*config, error) {
	config, err := bindConfiguration()
	if err != nil {
		log.Logger.DPanic(err.Error())
		panic(err)
	}
	cfg := getInstance()
	cfg.cfgConsumer = &kafka.ConfigMap{
		"bootstrap.servers":                   config.BootstrapServers,
		"group.id":                            config.GroupId,
		"security.protocol":                   config.SecurityProtocol,
		"ssl.certificate.location":            config.SslCertificateLocation,
		"message.max.bytes":                   config.MessageMaxBytes,
		"enable.ssl.certificate.verification": config.EnableSslCertificateVerification,
		"auto.offset.reset":                   config.AutoOffsetReset,
		"session.timeout.ms":                  config.SessionTimeoutMs,
		"partition.assignment.strategy":       config.PartitionAssignmentStrategy,
		"enable.auto.commit":                  true,
		"auto.commit.interval.ms":             500,
		"debug":                               config.ConsumerDebug,
	}
	cfg.cfgProducer = &kafka.ConfigMap{
		"bootstrap.servers":                   config.BootstrapServers,
		"security.protocol":                   config.SecurityProtocol,
		"ssl.certificate.location":            config.SslCertificateLocation,
		"message.max.bytes":                   config.MessageMaxBytes,
		"enable.ssl.certificate.verification": config.EnableSslCertificateVerification,
	}

	cfg.schemaRegistry = schemaregistry.NewConfig(config.SchemaRegistry)

	cfg.MaxRetry = config.MaxRetry

	return getInstance(), nil
}

func bindConfiguration() (*KafkaOption, error) {
	configuration := extension.GetConfiguration("enviroment.yaml")
	mapstructure.Decode(configuration["AMQStreams"], &configurations)

	err := validRequired()

	if err != nil {

		return nil, err
	}

	result := KafkaOption{
		BootstrapServers:                 getOrDefaultString(configurations, BootstrapServers, configurations[BootstrapServers]),
		SchemaRegistry:                   getOrDefaultString(configurations, SchemaRegistry, configurations[SchemaRegistry]),
		GroupId:                          getOrDefaultString(configurations, GroupId, ApplicationName),
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
	schemaRegistry := os.Getenv(SchemaRegistry)
	if schemaRegistry == "" && configurations[SchemaRegistry] == "" {
		return errors.New("the schemaRegistry is requiered")
	}
	return nil
}

func (c *config) ToConsumer(suscriber ISuscriber, event ISpecificRecord, topic []string) *config {

	subcription := ConsumerOptions{
		event:       event,
		topic:       topic,
		subscriptor: suscriber,
	}
	c.consumers = append(c.consumers, subcription)
	return c
}

func (c *config) Move(topic string) *config {
	c.consumers[len(c.consumers)-1].move = topic
	return c
}

func (c *config) ToProducer(event ISpecificRecord, topics []string) *config {
	appended := false

	for _, v := range c.producers {
		if v.ToPublish[event.Schema()] != nil {
			v.ToPublish[event.Schema()] = append(v.ToPublish[event.Schema()], topics...)
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
	return c
}

func (c *config) Build() {
	wg := new(sync.WaitGroup)
	for _, element := range c.consumers {
		wg.Add(1)
		event := element.event
		topic := element.topic
		go func() {
			for {
				err := c.consumer(event, topic, wg)
				if err != nil {
					log.SugarLogger.Errorln(err.Error())
				}
			}

		}()
	}
	wg.Wait()

}
