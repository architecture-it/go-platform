package AMQStream

type ISpecificRecord interface {
	Schema() string
	SchemaName() string
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}

type ISuscriber interface {
	Handler(event ISpecificRecord, metadata ConsumerMetadata)
}

type IPublisher interface {
	To(event ISpecificRecord, key string) error
}
