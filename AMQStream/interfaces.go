package AMQStream

type ISuscriber interface {
	Handler(event interface{}, metadata ConsumerMetadata)
}

type IPublisher interface {
	To(event interface{}, key string) error
}
