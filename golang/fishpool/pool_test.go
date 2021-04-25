package fishpool

import (
	"context"
	"testing"
	"time"
)

func TestPoolSizeAdjustment(t *testing.T) {
	pool := NewFunc(10, func(interface{}) interface{} {
		return "foo"
	})
	if exp, act := 10, len(pool.workers); exp != act {
		t.Errorf("Wrong size of pool: %v != %v", act, exp)
	}

	pool.SetSize(100)
	if exp, act := 100, pool.GetSize(); exp != act {
		t.Errorf("Wrong size of pool: %v != %v", act, exp)
	}

	pool.SetSize(5)
	if exp, act := 5, pool.GetSize(); exp != act {
		t.Errorf("Wrong size of pool: %v != %v", act, exp)
	}

	pool.Close()
	if exp, act := 0, pool.GetSize(); exp != act {
		t.Errorf("Wrong size of pool: %v != %v", act, exp)
	}
}


func TestFuncJob(t *testing.T) {
	pool := NewFunc(10, func(in interface{}) interface{} {
		intVal := in.(int)
		return intVal * 2
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		ret := pool.Process(10)
		if exp, act := 20, ret.(int); exp != act {
			t.Errorf("Wrong result: %v != %v", act, exp)
		}
	}
}

func TestFuncJobTimed(t *testing.T) {
	pool := NewFunc(10, func(in interface{}) interface{} {
		intVal := in.(int)
		return intVal * 2
	})
	defer pool.Close()

	for i := 0; i < 10; i++ {
		ret, err := pool.ProcessTimed(10, time.Millisecond)
		if err != nil {
			t.Fatalf("Failed to process: %v", err)
		}
		if exp, act := 20, ret.(int); exp != act {
			t.Errorf("Wrong result: %v != %v", act, exp)
		}
	}
}

func TestFuncJobCtx(t *testing.T) {
	t.Run("Completes when ctx not canceled", func(t *testing.T) {
		pool := NewFunc(10, func(in interface{}) interface{} {
			intVal := in.(int)
			return intVal * 2
		})
		defer pool.Close()

		for i := 0; i < 10; i++ {
			ret, err := pool.ProcessCtx(context.Background(), 10)
			if err != nil {
				t.Fatalf("Failed to process: %v", err)
			}
			if exp, act := 20, ret.(int); exp != act {
				t.Errorf("Wrong result: %v != %v", act, exp)
			}
		}
	})

	t.Run("Returns err when ctx canceled", func(t *testing.T) {
		pool := NewFunc(1, func(in interface{}) interface{} {
			inVal := in.(int)
			<-time.After(time.Millisecond)
			return inVal * 2
		})
		defer pool.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		_, act := pool.ProcessCtx(ctx, 10)
		if exp := context.DeadlineExceeded; exp != act {
			t.Errorf("Wrong error returned: %v != %v", act, exp)
		}
	})
}

// TODO  ADD TEST CASES