package AMQStream

func To(event ISpecificRecord, key string) error {

	return getInstance().to(event, key)
}
