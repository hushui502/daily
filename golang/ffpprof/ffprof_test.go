package ffpprof

import (
	"strconv"
	"testing"
	"time"
)

func TestBlockProfile_Capture(t *testing.T) {
	Capture(CPUProfile{
		Duration: 30 * time.Second,
	})

	for i := 0; i < 5000; i++ {
		generateID(5, 1000)
	}
}

func generateID(duration int, usage int) {
	for j := 0; j < duration; j++ {
		go func() {
			for i := 0; i < usage*10000; i++ {
				str := "str" + strconv.Itoa(i)
				str = str + "a"
			}
		}()
		time.Sleep(1 * time.Second)
	}
}