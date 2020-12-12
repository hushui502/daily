package eventbus

import (
	"testing"
	"time"
)

func sub1(args ...interface{}) {
	time.Sleep(1 * time.Microsecond)
	//fmt.Printf("sub1, %s %s\n", msg1, msg2)
}

func sub2(args ...interface{}) {
	//fmt.Printf("sub2, %s %s\n", msg1, msg2)
}

func TestAsyncEventBus_Publish(t *testing.T) {
	bus := NewAsyncEventBus()
	bus.Subscribe("topic1", sub1)
	bus.Subscribe("topic1", sub2)

	bus.Publish("topic1", "aaa", "bbb", "ccc")
}
