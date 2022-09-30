package goncu_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/rodriez/goncu/v2"
)

func TestPool_Run_Without_DO(t *testing.T) {
	err := goncu.Pool(10, 50).Run(nil)

	if err == nil {
		t.Error("expected error but get nil")
	}
}

func TestPool_Run_With_invalid_task_amount(t *testing.T) {
	err := goncu.Pool(10, -10).Run(func(n int) {
		t.Error("this function should not be executed")
	})

	if err == nil {
		t.Error("expected error but get nil")
	}
}

func TestPool_Run_With_invalid_workers_amount(t *testing.T) {
	hits := int32(0)
	err := goncu.Pool(-1, 10).Run(func(n int) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if hits != 10 {
		t.Errorf("expected 10 hits but get %d", hits)
	}
}

func TestPool_Run_With_more_workers_than_task(t *testing.T) {
	hits := int32(0)
	err := goncu.Pool(100, 10).Run(func(n int) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if hits != 10 {
		t.Errorf("expected 10 hits but get %d", hits)
	}
}

func TestPool_Run_Ok(t *testing.T) {
	hits := int32(0)
	err := goncu.Pool(10, 50).Run(func(n int) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if hits != 50 {
		t.Errorf("expected a 50 hits but get %d", hits)
	}
}

func TestIteratorPool_With_invalid_workers_amount(t *testing.T) {
	iterator := goncu.Iterator(100, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	hits := int32(0)
	err := goncu.IteratorPool(-1, iterator).Run(func(s string) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if hits != 100 {
		t.Errorf("expected a 100 hits but get %d", hits)
	}
}

func TestIteratorPool_Run_Ok(t *testing.T) {
	iterator := goncu.Iterator(100, func(i int) string {
		return fmt.Sprintf("%d", i)
	})

	buffer := make(chan string, 100)
	defer close(buffer)
	err := goncu.IteratorPool(10, iterator).Run(func(s string) {
		buffer <- s
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if len(buffer) != 100 {
		t.Errorf("expected a 100 hits but get %d", len(buffer))
	}
}

func TestSlicePool_With_invalid_workers_amount(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	hits := int32(0)
	err := goncu.SlicePool(-1, data).Run(func(s string) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if int(hits) != len(data) {
		t.Errorf("expected a %d hits but get %d", len(data), hits)
	}
}

func TestSlicePool_Run_Ok(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	hits := int32(0)
	err := goncu.SlicePool(3, data).Run(func(s string) {
		atomic.AddInt32(&hits, 1)
	})

	if err != nil {
		t.Errorf("no error was expected but get %s", err)
	}

	if int(hits) != len(data) {
		t.Errorf("expected a %d hits but get %d", len(data), hits)
	}
}
