package orderedmap

import (
	"testing"
)

func TestNewOrderedMapEmpty(t *testing.T) {
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
}

func TestNewOrderedMapAddElements(t *testing.T) {
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
}

func TestNewOrderedMapDifferentTypes(t *testing.T) {
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
}
func TestPutAndGet(t *testing.T) {
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
}

func TestDelete(t *testing.T) {
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
}

func TestPutAndGetWithDifferentTypes(t *testing.T) {
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
}

const numberOfElements = 1000
