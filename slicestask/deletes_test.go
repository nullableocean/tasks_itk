package slicestask

import (
	"reflect"
	"testing"
)

func TestDeletes_Work(t *testing.T) {
	t.Run("RemoveUnordered", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		result := RemoveUnordered(s, 2)
		expected := []int{1, 2, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveUnordered got %v, want %v", result, expected)
		}

		// out of bounds
		result = RemoveUnordered(s, 10)
		if !reflect.DeepEqual(result, s) {
			t.Errorf("RemoveUnordered out range should return original")
		}
	})

	t.Run("RemoveOrdered", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}
		result := RemoveOrdered(s, 2)
		expected := []int{1, 2, 0, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveOrdered got %v, want %v", result, expected)
		}
	})

	t.Run("RemoveAllByValue", func(t *testing.T) {
		s := []int{1, 2, 2, 3, 2, 4, 2}
		result := RemoveAllByValue(s, 2)
		expected := []int{1, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveAllByValue got %v, want %v", result, expected)
		}

		// no such value
		result = RemoveAllByValue(s, 99)
		if !reflect.DeepEqual(result, s) {
			t.Errorf("RemoveAllByValue non-existent should return original")
		}
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		s := []int{1, 2, 2, 3, 1, 4, 2, 5}
		result := RemoveDuplicates(s)
		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveDuplicates got %v, want %v", result, expected)
		}

		// all unique
		s2 := []int{1, 2, 3}
		result2 := RemoveDuplicates(s2)
		if !reflect.DeepEqual(result2, s2) {
			t.Errorf("RemoveDuplicates all unique should return original")
		}
	})

	t.Run("RemoveIf", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5, 6}
		result := RemoveIf(s, func(x int) bool { return x%2 == 0 })
		expected := []int{1, 3, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("RemoveIf = %v, want %v", result, expected)
		}

		// none removed
		result = RemoveIf(s, func(x int) bool { return false })
		if !reflect.DeepEqual(result, s) {
			t.Errorf("RemoveIf none should return original")
		}
	})

	t.Run("RemoveOrderedWithNil", func(t *testing.T) {
		val1, val2, val3 := 1, 2, 3
		s := []*int{&val1, &val2, &val3}
		result := RemoveOrderedWithNil(s, 1)

		if result[1] != nil {
			t.Errorf("RemoveOrderedWithNil should nil the element")
		}
		if len(result) != 3 {
			t.Errorf("RemoveOrderedWithNil length = %d, want 3", len(result))
		}
	})

	t.Run("ShrinkCapacity", func(t *testing.T) {
		// no shrink needed
		s := make([]int, 5, 10)
		result := ShrinkCapacity(s)
		if cap(result) != 10 {
			t.Errorf("ShrinkCapacity should not shrink when cap < 2*len")
		}

		// shrink needed
		s2 := make([]int, 5, 20)
		result2 := ShrinkCapacity(s2)
		if cap(result2) != 5 {
			t.Errorf("ShrinkCapacity cap = %d, want 5", cap(result2))
		}
	})
}
