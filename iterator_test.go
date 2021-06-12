package goncu_test

import (
	"testing"

	"github.com/rodriez/goncu"
)

func TestIterator_Ok(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	iterator := goncu.NewIterator(len(slice)).
		Producer(func(i int) interface{} {
			return slice[i]
		})

	newSlice := []string{}
	for e := range iterator.Each() {
		newSlice = append(newSlice, e.(string))
	}

	if len(newSlice) < len(slice) {
		t.Error("there is less elements than expected")
	}

	if len(newSlice) > len(slice) {
		t.Error("there is more elements than expected")
	}
}

func TestIterator_With_Empty_Producer(t *testing.T) {
	iterator := goncu.NewIterator(10)

	newSlice := []string{}
	for e := range iterator.Each() {
		newSlice = append(newSlice, e.(string))
	}

	if len(newSlice) != 0 {
		t.Error("there is more elements than expected")
	}
}
