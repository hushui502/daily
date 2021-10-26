package go_mq

type Consumer interface {
	Consume(delivery Delivery)
}

type ConsumerFunc func(Delivery)

func (consumerFunc ConsumerFunc) Consume(delivery Delivery) {
	consumerFunc(delivery)
}

type BatchConsumer interface {
	Consume(batch Deliveries)
}
