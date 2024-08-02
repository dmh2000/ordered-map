package orderedmap

import (
	"testing"
)

func TestGetPutDelete(t *testing.T) {
	t.Run("Put and Get", func(t *testing.T) {
		om := NewOrderedMap[int, string]()
		
		// Put some key-value pairs
		om.Put(1, "one")
		om.Put(2, "two")
		om.Put(3, "three")

		// Test Get
		value, found := om.Get(2)
		if !found {
			t.Error("Key 2 not found in map")
		}
		if value != "two" {
			t.Errorf("Expected value 'two' for key 2, got '%s'", value)
		}

		// Test Get for non-existent key
		_, found = om.Get(4)
		if found {
			t.Error("Found non-existent key 4 in map")
		}

		// Test overwriting a value
		om.Put(2, "TWO")
		value, _ = om.Get(2)
		if value != "TWO" {
			t.Errorf("Expected updated value 'TWO' for key 2, got '%s'", value)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		om := NewOrderedMap[int, string]()
		
		// Put some key-value pairs
		om.Put(1, "one")
		om.Put(2, "two")
		om.Put(3, "three")

		// Delete existing key
		om.Delete(2)
		_, found := om.Get(2)
		if found {
			t.Error("Key 2 still found in map after deletion")
		}

		// Check size after deletion
		if om.Size() != 2 {
			t.Errorf("Expected size 2 after deletion, got %d", om.Size())
		}

		// Delete non-existent key (should not panic or affect size)
		om.Delete(4)
		if om.Size() != 2 {
			t.Errorf("Expected size to remain 2 after deleting non-existent key, got %d", om.Size())
		}

		// Delete remaining keys
		om.Delete(1)
		om.Delete(3)
		if !om.IsEmpty() {
			t.Error("Map should be empty after deleting all keys")
		}
	})

	t.Run("Put and Get with different types", func(t *testing.T) {
		omString := NewOrderedMap[string, int]()
		omString.Put("one", 1)
		omString.Put("two", 2)

		value, found := omString.Get("one")
		if !found {
			t.Error("Key 'one' not found in string map")
		}
		if value != 1 {
			t.Errorf("Expected value 1 for key 'one', got %d", value)
		}

		omFloat := NewOrderedMap[float64, bool]()
		omFloat.Put(1.1, true)
		omFloat.Put(2.2, false)

		value2, found := omFloat.Get(2.2)
		if !found {
			t.Error("Key 2.2 not found in float map")
		}
		if value2 != false {
			t.Errorf("Expected value false for key 2.2, got %v", value2)
		}
	})
}
