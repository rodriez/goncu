package goncu_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rodriez/goncu"
)

func TestPool_Run_Without_DO(t *testing.T) {
	_, err := goncu.NewPool(10).OnSuccess(func(item interface{}) {
		fmt.Println(item)
	}).OnError(func(e error) {
		fmt.Println(e)
	}).Run(50)

	if err == nil {
		t.Error("expected error but get nil")
	}
}

func TestPool_Run_With_invalid_task_amount(t *testing.T) {
	_, err := goncu.NewPool(10).DO(func(n int) (interface{}, error) {
		if n%2 != 0 {
			return nil, errors.New("oops")
		}

		time.Sleep(300 * time.Millisecond)
		return n + 1, nil
	}).OnSuccess(func(item interface{}) {
		fmt.Println(item)
	}).OnError(func(e error) {
		fmt.Println(e)
	}).Run(-10)

	if err == nil {
		t.Error("expected error but get nil")
	}
}

func TestPool_Run_With_invalid_workers_amount(t *testing.T) {
	response, err := goncu.NewPool(-1).DO(func(n int) (interface{}, error) {
		return n + 1, nil
	}).OnSuccess(func(item interface{}) {
		fmt.Println(item)
	}).OnError(func(e error) {
		fmt.Println(e)
	}).Run(10)

	if err != nil {
		t.Errorf("expected nil but get %s", err)
	}

	if response == nil {
		t.Error("expected a response but get nil")
	} else if len(response.Errors) != 0 {
		t.Errorf("expected a no errors but get %d", len(response.Errors))
	} else if len(response.Hits) != 10 {
		t.Errorf("expected a 10 hits but get %d", len(response.Hits))
	}
}

func TestPool_Run_Ok(t *testing.T) {
	response, err := goncu.NewPool(10).DO(func(n int) (interface{}, error) {
		if n%2 != 0 {
			return nil, errors.New("oops")
		}

		time.Sleep(300 * time.Millisecond)
		return n + 1, nil
	}).OnSuccess(func(item interface{}) {
		fmt.Println(item)
	}).OnError(func(e error) {
		fmt.Println(e)
	}).Run(50)

	if err != nil {
		t.Errorf("expected nil but get %s", err)
	}

	if response == nil {
		t.Error("expected a response but get nil")
	} else if len(response.Errors) != 25 {
		t.Errorf("expected a 25 errors but get %d", len(response.Errors))
	} else if len(response.Hits) != 25 {
		t.Errorf("expected a 25 hits but get %d", len(response.Hits))
	}
}
