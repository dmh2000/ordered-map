package orderedmap

import (
	"testing"
)

func TestNewOrderedMap(t *testing.T) {
	t.Run("Create empty map", func(t *testing.T) {
		om := NewOrderedMap[int, string]()
		if om == nil {
			t.Error("NewOrderedMap returned nil")
		}
		if !om.IsEmpty() {
			t.Error("Newly created map is not empty")
		}
		if om.Size() != 0 {
			t.Errorf("Expected size 0, got %d", om.Size())
		}
	})

	t.Run("Create and add elements", func(t *testing.T) {
		om := NewOrderedMap[int, string]()
		om.Put(1, "one")
		om.Put(2, "two")
		om.Put(3, "three")

		if om.IsEmpty() {
			t.Error("Map should not be empty after adding elements")
		}
		if om.Size() != 3 {
			t.Errorf("Expected size 3, got %d", om.Size())
		}

		value, found := om.Get(2)
		if !found {
			t.Error("Key 2 not found in map")
		}
		if value != "two" {
			t.Errorf("Expected value 'two' for key 2, got '%s'", value)
		}
	})

	t.Run("Create with different types", func(t *testing.T) {
		omString := NewOrderedMap[string, int]()
		omString.Put("one", 1)
		omString.Put("two", 2)

		if omString.Size() != 2 {
			t.Errorf("Expected size 2 for string map, got %d", omString.Size())
		}

		omFloat := NewOrderedMap[float64, bool]()
		omFloat.Put(1.1, true)
		omFloat.Put(2.2, false)

		if omFloat.Size() != 2 {
			t.Errorf("Expected size 2 for float map, got %d", omFloat.Size())
		}
	})
}
