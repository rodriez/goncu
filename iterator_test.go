package goncu_test

import (
	"testing"

	"github.com/rodriez/goncu"
)

func TestIterator_Ok(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	iterator := goncu.Iterator(12, func(i int) string {
		return slice[i]
	})

	newSlice := []string{}
	iterator.ForEach(func(i int, s string) {
		newSlice = append(newSlice, s)
	})

	if len(newSlice) < len(slice) {
		t.Error("there is less elements than expected")
	}

	if len(newSlice) > len(slice) {
		t.Error("there is more elements than expected")
	}
}

func TestIterator_With_Empty_Producer(t *testing.T) {
	iterator := goncu.Iterator[string](10, nil)

	newSlice := []string{}
	iterator.ForEach(func(i int, s string) {
		newSlice = append(newSlice, s)
	})

	if len(newSlice) != 0 {
		t.Error("there is more elements than expected")
	}
}

func TestSliceIterator_Ok(t *testing.T) {
	slice := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

	iterator := goncu.SliceIterator(slice)

	newSlice := iterator.ToArray()

	if len(newSlice) < len(slice) {
		t.Error("there is less elements than expected")
	}

	if len(newSlice) > len(slice) {
		t.Error("there is more elements than expected")
	}
}
