package orderedmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// TestNewOrderedMapEmpty tests the behavior of NewOrderedMap when an empty map is created.
//
// This function creates a new OrderedMap with integer keys and string values.
// It checks that the map is not nil, is empty, and has a size of 0.
// If any of these conditions are not met, an error is reported.
// This function takes a *testing.T object as a parameter and does not return anything.
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

// TestNewOrderedMapAddElements tests the behavior of NewOrderedMap when elements are added to it.
//
// This function creates a new OrderedMap with integer keys and string values.
// It adds three key-value pairs to the map using the Put method.
// It then checks that the map is not empty and has a size of 3.
// It retrieves the value associated with the key 2 using the Get method and checks that it is equal to "two".
// If any of these conditions are not met, an error is reported.
// This function takes a *testing.T object as a parameter and does not return anything.
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

// TestNewOrderedMapDifferentTypes tests the creation of two different types of OrderedMap
// and checks if the size of each map is 2.
//
// t *testing.T - the testing object.
// No return value.
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

// TestPutAndGet tests the Put and Get methods of the OrderedMap.
//
// It creates a new OrderedMap and puts some key-value pairs into it.
// Then it tests the Get method by retrieving a value for an existing key
// and checking if the value is correct. It also tests the Get method for a
// non-existent key and checks if it returns false. Finally, it tests the
// overwriting of a value by putting a new value for an existing key and
// checking if the updated value is correct.
//
// Parameters:
// - t: The testing.T object for running the test.
//
// Returns: None.
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

// TestDelete is a function to test the Delete method of the OrderedMap.
//
// - t: The testing.T object for running the test.
// Returns: None.
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

// TestPutAndGetWithDifferentTypes is a function to test the Put and Get methods of the OrderedMap with different types.
//
// This function creates two OrderedMaps, one with string keys and integer values, and another with float64 keys and boolean values.
// It puts some key-value pairs in each map using the Put method.
// It then retrieves the value associated with a key in each map using the Get method and checks that it is correct.
// If any of the retrieved values are not correct, an error is reported.
//
// Parameters:
// - t: The testing.T object for running the test.
//
// Returns: None.
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

const numberOfElements = 5000000

// uniqueValues generates a map of unique random key/value pairs.
//
// Parameters:
// - n: the number of key/value pairs to generate.
//
// Returns:
// - map[string]int: a map containing unique random key/value pairs.
func uniqueValues(n int) map[string]int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	m := make(map[string]int)
	for i := 0; i < n; i++ {
		v := rng.Int()
		k := strconv.Itoa(rng.Int())
		m[k] = v
	}
	return m
}

// TestRandomAddGetDelete tests the functionality of adding, getting, and deleting
// elements from an OrderedMap. It creates a list of unique random key/value pairs,
// adds them to the OrderedMap, verifies the size of the OrderedMap, gets and verifies
// the values of the elements, and deletes the elements from the OrderedMap.
//
// Parameters:
// - t: the testing.T object used for reporting test failures.
//
// Return type: None.
func TestRandomAddGetDelete(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// create a list of unique random key/value pairs
	m := uniqueValues(numberOfElements)
	elements := len(m)

	// Add random elements
	for k, v := range m {
		om.Put(k, v)
	}

	if om.Size() != elements {
		t.Errorf("Expected size %d, got %d", numberOfElements, om.Size())
	}

	// Get and verify random elements
	for k, v := range m {
		value, found := om.Get(k)
		if found {
			if value != v {
				t.Errorf("keys[i] = %s values[i] = %d, got %d", k, value, v)
			}
		}
	}

	// Delete random elements
	expectedSize := elements
	for k := range m {
		om.Delete(k)
		if om.Size() != expectedSize-1 {
			t.Errorf("Expected size %d after deletion, got %d", expectedSize-1, om.Size())
		}
		expectedSize--
	}
}

// TestIterateOverKeys tests the functionality of iterating over keys in an OrderedMap.
//
// It creates a new OrderedMap and adds a list of unique random key/value pairs to it.
// Then it verifies that the size of the OrderedMap matches the number of elements added.
// Next, it iterates over the keys and checks that all keys are found and their corresponding values match the original map.
// It also checks that the keys are in the correct order.
// Finally, it deletes random elements from the OrderedMap and verifies that the size decreases correctly.
//
// Parameters:
// - t: A testing.T object for reporting test failures.
//
// Return type: None.
func TestIterateOverKeys(t *testing.T) {
	om := NewOrderedMap[string, int]()

	// create a list of unique random key/value pairs
	m := uniqueValues(numberOfElements)
	elements := len(m)

	// Add random elements
	for k, v := range m {
		om.Put(k, v)
	}

	if om.Size() != elements {
		t.Errorf("Expected size %d, got %d", numberOfElements, om.Size())
	}

	// Get and verify random elements
	keys := om.Keys()
	for _, k := range keys {
		value, found := om.Get(k)
		// check all keys are found
		if !found {
			t.Errorf("Key %s not found in map", k)
		}
		// check value against map
		if found {
			if value != m[k] {
				t.Errorf("keys[i] = %s values[i] = %d, got %d", k, value, m[k])
			}
		}
	}

	// check keys are in order
	keys = om.Keys()
	kn := keys[0]
	for i := 1; i < len(keys); i++ {
		if kn >= keys[i] {
			t.Errorf("Keys not in order: %s, %s", kn, keys[i])
		}
		// t.Log(kn, keys[i])
		kn = keys[i]
	}

	// Delete random elements
	expectedSize := elements
	for k := range m {
		om.Delete(k)
		if om.Size() != expectedSize-1 {
			t.Errorf("Expected size %d after deletion, got %d", expectedSize-1, om.Size())
		}
		expectedSize--
	}
}

// TestIterateOverMap tests the iteration over a map created from unique random key/value pairs.
// used as a benchmark for map vs OrderedMap
//
// Parameters:
// - t: the testing.T object used for reporting test failures.
// Return type: None.
func TestIterateOverMap(t *testing.T) {
	// create a list of unique random key/value pairs
	m := uniqueValues(numberOfElements)

	// create a new maps from m
	p := make(map[string]int)
	for k, v := range m {
		p[k] = v
	}

	// Get and verify random elements
	for k, v := range p {
		value, found := p[k]
		if found {
			if value != v {
				t.Errorf("keys[i] = %s values[i] = %d, got %d", k, value, v)
			}
		}
	}
}

// TestKeysInRangeBFS tests the KeysInRangeBFS function of the OrderedMap type.
//
// This function creates a new OrderedMap with 32 key-value pairs, where the keys are integers
// and the values are strings representing the same integers. It then prints the keys in the
// order they were inserted. Next, it calls the KeysInRangeBFS function to get the keys in
// breadth-first search order and checks if the length of the returned slice is 32. If not,
// it reports an error using the t.Errorf function. It then prints the keys again. Finally,
// it prints the keys in a tree-like structure, where each level of the tree is printed on a
// new line.
//
// Parameters:
// - t: a testing.T object used for reporting errors and logging information.
//
// Return type: None.
func TestKeysInRangeBFS(t *testing.T) {
	om := NewOrderedMap[int, string]()
	for i := 0; i < 32; i++ {
		om.Put(i, strconv.Itoa(i))
	}

	// print in order
	keys := om.Keys()
	fmt.Println(keys)

	// test keysInRangeBFS
	keys = om.KeysInRangeBFS()
	if len(keys) != 32 {
		t.Errorf("Expected 32 keys, got %d", len(keys))
	}

	fmt.Println(keys)

	level := 1
	step := 0
	for _, k := range keys {
		fmt.Printf("%v ", k)
		step++
		if step == level {
			fmt.Println()
			level *= 2
			step = 0
		}
	}
	fmt.Println()
}

// FuzzOrderedMap is a fuzz testing function that takes a *testing.F as input.
//
// It creates a new OrderedMap[int32, string] and adds some key-value pairs to it.
// Then, it adds random elements to the map using the provided fuzzing function.
// For each element, it checks if the key exists in the map, if the value matches,
// deletes the key from the map, checks if the key is not found, and checks if the
// size of the map is 0 after deletion.
//
// Parameters:
// - f: a pointer to a *testing.F object.
//
// Return type: None.
func FuzzOrderedMap(f *testing.F) {
	om := NewOrderedMap[int32, string]()

	// Add some key-value pairs
	f.Add(int32(1), "abc")
	f.Add(int32(2), "twssso")
	f.Add(int32(3), "thtrxree")

	// Add random elements
	f.Fuzz(func(t *testing.T, key int32, value string) {
		om.Put(key, value)
		val, found := om.Get(key)
		if !found {
			t.Errorf("Key %d not found in map", key)
		}
		if val != value {
			t.Errorf("Expected value %s for key %d, got %s", value, key, val)
		}
		om.Delete(key)
		_, found = om.Get(key)
		if found {
			t.Errorf("Key %d found in map after deletion", key)
		}
		sz := om.Size()
		if sz != 0 {
			t.Errorf("Expected size 0, got %d", sz)
		}
	})
}
