package ffcmp

import (
	"runtime"
	"testing"
)

func TestBasics(t *testing.T) {
	foo := Get("foo")
	bar := Get("bar")
	foo2 := Get("foo")
	bar2 := Get("bar")

	if foo.Get() != foo2.Get() {
		t.Error("foo values differ")
	}

	if bar.Get() != bar2.Get() {
		t.Error("bar values differ")
	}

	if foo != foo2 {
		t.Error("foo pointer differ")
	}

	if bar != bar2 {
		t.Error("bar pointer differ")
	}

	const gcTries = 5000
	for try := 0; try < gcTries; try++ {
		runtime.GC()
		n := mapLen()
		if n == 0 {
			break
		}

		if try == gcTries-1 {
			t.Errorf("map len = %d after (%d GC tries); want 0", gcTries, try)
		}
	}

}

func mapLen() int {
	mu.Lock()
	defer mu.Unlock()

	return len(valMap)
}
