package ringbuffer_pro

import (
	"sync"
	"testing"
)

func BenchmarkConsumerSequentialReadWriteLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerSequentialReadWrite(10000000)
	}
}

func BenchmarkChannelsSequentialReadWriteLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsSequentialReadWrite(10000000)
	}
}

func BenchmarkConsumerSequentialReadWriteMedium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerSequentialReadWrite(100000)
	}
}

func BenchmarkChannelsSequentialReadWriteMedium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsSequentialReadWrite(100000)
	}
}

func BenchmarkConsumerSequentialReadWriteSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerSequentialReadWrite(1000)
	}
}

func BenchmarkChannelsSequentialReadWriteSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsSequentialReadWrite(1000)
	}
}

func ConsumerSequentialReadWrite(n int) {

	var buffer = CreateBuffer[int](BufferSize, 10)
	consumer, _ := buffer.CreateConsumer()

	for i := 0; i < n; i++ {
		buffer.Write(i)
		consumer.Get()
	}
}

func ChannelsSequentialReadWrite(n int) {

	var buffer = make(chan int, BufferSize)

	for i := 0; i < n; i++ {
		buffer <- i
		<-buffer
	}
}

/*
General Benchmark to compare concurrent reading from channels vrs the ring buffer.
Note there is heavy over head for syncing the routines in both and is not accurate beyond a general comparison.
*/
func BenchmarkConsumerConcurrentReadWriteLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerConcurrentReadWrite(10000000)
	}
}

func BenchmarkChannelsConcurrentReadWriteLarge(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsConcurrentReadWrite(10000000)
	}
}

func BenchmarkConsumerConcurrentReadWriteMedium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerConcurrentReadWrite(100000)
	}
}

func BenchmarkChannelsConcurrentReadWriteMedium(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsConcurrentReadWrite(100000)
	}
}
func BenchmarkConsumerConcurrentReadWriteSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ConsumerConcurrentReadWrite(1000)
	}
}

func BenchmarkChannelsConcurrentReadWriteSmall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ChannelsConcurrentReadWrite(1000)
	}
}

func ConsumerConcurrentReadWrite(n int) {

	var buffer = CreateBuffer[int](BufferSize, 10)

	var wg sync.WaitGroup
	messages := []int{}

	for i := 0; i < n; i++ {
		messages = append(messages, i)
	}

	consumer, _ := buffer.CreateConsumer()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, value := range messages {
			buffer.Write(value)
		}
	}()

	i := -1

	wg.Add(1)
	go func() {

		defer wg.Done()
		for _, _ = range messages {
			j := consumer.Get()
			if j != i+1 {
				panic("data is inconsistent")
			}
			i = j
		}
	}()
	wg.Wait()
}

func ChannelsConcurrentReadWrite(n int) {

	var wg sync.WaitGroup
	messages := []int{}

	for i := 0; i < n; i++ {
		messages = append(messages, i)
	}

	var buffer = make(chan int, BufferSize)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, value := range messages {
			buffer <- value
		}
	}()

	i := -1

	wg.Add(1)
	go func() {

		defer wg.Done()
		for _, _ = range messages {
			j := <-buffer
			if j != i+1 {
				panic("data is inconsistent")
			}
			i = j
		}
	}()
	wg.Wait()
}

