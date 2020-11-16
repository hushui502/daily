package pattern

import "fmt"

type Consumer interface {
	update()
}

type Consumer1 struct {

}

func (c1 *Consumer1) update() {
	fmt.Println("c1 received message")
}

type Consumer2 struct {

}

func (c2 *Consumer2) update() {
	fmt.Println("c2 received message")
}

type NewOffice struct {
	consumers []Consumer
}

func (n *NewOffice) addConsumer(consumer Consumer) {
	n.consumers = append(n.consumers, consumer)
}

func (n *NewOffice) newspaperCome() {
	n.notifyAllConsumer()
}

func (n *NewOffice) notifyAllConsumer() {
	for _, consumer := range n.consumers {
		consumer.update()
	}
}

