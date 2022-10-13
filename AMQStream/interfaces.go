package AMQStream

import (
	"io"

	"github.com/actgardner/gogen-avro/v10/vm/types"
)

type ISpecificRecord interface {
	types.Field
	Serialize(w io.Writer) error
	Schema() string
	SchemaName() string
}

type ISuscriber interface {
	Handler(event interface{}, metadata ConsumerMetadata) error
}

type IPublisher interface {
	To(event interface{}, key string) error
}
