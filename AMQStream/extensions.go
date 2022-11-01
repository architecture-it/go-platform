package AMQStream

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

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
	cfg.cfgAdmin = &kafka.ConfigMap{
		"bootstrap.servers":                   config.BootstrapServers,
		"group.id":                            config.GroupId,
		"security.protocol":                   config.SecurityProtocol,
		"ssl.certificate.location":            config.SslCertificateLocation,
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

func (c *config) CreateOrUpdateTopics(numParts int, topics []string) *config {
	replicationFactor := 3

	// Create a new AdminClient.
	// AdminClient can also be instantiated using an existing
	// Producer or Consumer instance, see NewAdminClientFromProducer and
	// NewAdminClientFromConsumer.
	adminClient, err := kafka.NewAdminClient(c.cfgAdmin)
	if err != nil {
		log.SugarLogger.Infoln("Failed to create Admin client: %s\n", err)
		os.Exit(1)
	}

	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("60s")
	if err != nil {
		defer time.ParseDuration("60s")
	}
	for _, topic := range topics {
		results, err := adminClient.CreateTopics(
			ctx,
			// Multiple topics can be created simultaneously
			// by providing more TopicSpecification structs here.
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     numParts,
				ReplicationFactor: replicationFactor}},
			// Admin options
			kafka.SetAdminOperationTimeout(maxDur))
		if err != nil {
			log.SugarLogger.Infoln("Failed to create topic: %v\n", err)
		}
		// Print results
		for _, result := range results {
			if result.Error.Code() == kafka.ErrTopicAlreadyExists {
				updateTopic(adminClient, ctx, topic, numParts, maxDur)
			} else {
				log.SugarLogger.Infoln("%s\n", result)
			}
		}
	}
	adminClient.Close()
	return c
}

func updateTopic(adminClient *kafka.AdminClient, ctx context.Context, topic string, numParts int, maxDur time.Duration) {
	resultUpdate, err := adminClient.CreatePartitions(ctx, []kafka.PartitionsSpecification{{
		Topic:      topic,
		IncreaseTo: numParts,
	}}, kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		log.SugarLogger.Infoln("Failed to Update topic: %v\n", err)
	}
	for _, result := range resultUpdate {
		if result.Error.Code() == kafka.ErrNoError || result.Error.Code() == kafka.ErrInvalidPartitions {

		} else {
			log.SugarLogger.Infoln("%s\n", result)
		}
	}
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
