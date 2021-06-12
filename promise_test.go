package goncu_test

import (
	"testing"
	"time"

	"github.com/rodriez/goncu"
)

func TestPromise(t *testing.T) {
	promise := goncu.NewPromise().
		Start(func() (interface{}, error) {
			return "data", nil
		})

	time.Sleep(100 * time.Millisecond)

	resp, err := promise.Done()

	if err != nil {
		t.Errorf("expected no error but get %s", err)
	}

	if str := resp.(string); str != "data" {
		t.Errorf("expected data but get %s", str)
	}
}
