package goncu_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/rodriez/goncu/v2"
)

func TestAsyncPromise(t *testing.T) {
	thenCalled := false

	goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)

		return "data", nil
	}).Then(func(str string) {
		thenCalled = true
		if str != "data" {
			t.Errorf("expected data but get %s", str)
		}
	}).Catch(func(err error) {
		t.Errorf("expected no error but get %s", err)
	})

	time.Sleep(100 * time.Millisecond)

	if !thenCalled {
		t.Errorf("Promise.Then was not executed")
	}
}

func TestAsyncPromiseWithError(t *testing.T) {
	catchCalled := false

	goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)

		return "", fmt.Errorf("oops")
	}).Then(func(str string) {
		t.Errorf("no data was expected but get %s", str)
	}).Catch(func(err error) {
		catchCalled = true
		if err == nil {
			t.Error("an error was expected but get nil")
		}
	})

	time.Sleep(100 * time.Millisecond)

	if !catchCalled {
		t.Errorf("Promise.Catch was not executed")
	}
}

func TestSyncPromise(t *testing.T) {
	promise := goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "data", nil
	})

	time.Sleep(100 * time.Millisecond)

	str, err := promise.Done()

	if str != "data" {
		t.Errorf("expected data but get %s", str)
	}

	if err != nil {
		t.Errorf("expected no error but get %s", err)
	}
}

func TestSyncPromiseWithError(t *testing.T) {
	promise := goncu.NewPromise(func() (string, error) {
		time.Sleep(50 * time.Millisecond)
		return "", fmt.Errorf("oops")
	})

	time.Sleep(100 * time.Millisecond)

	str, err := promise.Done()

	if str != "" {
		t.Errorf("no data was expected but get %s", str)
	}

	if err == nil {
		t.Error("an error was expected but get nil")
	}
}
