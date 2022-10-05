package AMQStream

func To(event interface{}, key string) error {

	return getInstance().to(event, key)
}
