package go_mq

type Deliveries []Delivery

func (ds Deliveries) Payload() []string {
	payloads := make([]string, len(ds))
	for i, delivery := range ds {
		payloads[i] = delivery.Payload()
	}

	return payloads
}

func (ds Deliveries) Ack() (errMap map[int]error) {
	return ds.each(Delivery.Ack)
}

func (ds Deliveries) Reject() (errMap map[int]error) {
	return ds.each(Delivery.Reject)
}

func (ds Deliveries) Push() (errMap map[int]error) {
	return ds.each(Delivery.Push)
}

func (ds Deliveries) each(f func(Delivery) error) (errMap map[int]error) {
	for i, delivery := range ds {
		if err := f(delivery); err != nil {
			if errMap == nil {
				errMap = map[int]error{}
			}
			errMap[i] = err
		}
	}

	return errMap
}
