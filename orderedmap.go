package orderedmap

import "golang.org/x/exp/constraints"

type color bool

const (
	RED   color = true
	BLACK color = false
)

type node[K constraints.Ordered, V any] struct {
	key         K
	val         V
	left, right *node[K, V]
	color       color
	size        int
}

type OrderedMap[K constraints.Ordered, V any] struct {
	root *node[K, V]
}

// NewOrderedMap creates and returns a new empty OrderedMap.
func NewOrderedMap[K constraints.Ordered, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

// Get retrieves the value associated with the given key.
func (t *OrderedMap[K, V]) Get(key K) (V, bool) {
	return t.get(t.root, key)
}

// Put inserts a key-value pair into the OrderedMap.
// If the key already exists, its value is updated.
func (t *OrderedMap[K, V]) Put(key K, val V) {
	t.root = t.put(t.root, key, val)
	t.root.color = BLACK
}

// Contains checks if the given key exists in the OrderedMap.
func (t *OrderedMap[K, V]) Contains(key K) bool {
	_, found := t.Get(key)
	return found
}

// Delete removes the key-value pair with the given key from the OrderedMap.
// If the key doesn't exist, this operation does nothing.
func (t *OrderedMap[K, V]) Delete(key K) {
	if !t.Contains(key) {
		return
	}

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}

	t.root = t.delete(t.root, key)
	if !t.IsEmpty() {
		t.root.color = BLACK
	}
}

// Keys returns a slice containing all keys in the OrderedMap in sorted order.
func (t *OrderedMap[K, V]) Keys() []K {
	if t.IsEmpty() {
		return []K{}
	}
	min, _ := t.Min()
	max, _ := t.Max()
	return t.KeysInRange(min, max)
}

// Size returns the number of key-value pairs in the OrderedMap.
func (t *OrderedMap[K, V]) Size() int {
	return t.size(t.root)
}

// IsEmpty returns true if the OrderedMap contains no elements, false otherwise.
func (t *OrderedMap[K, V]) IsEmpty() bool {
	return t.root == nil
}

// Min returns the smallest key in the OrderedMap and a boolean indicating success.
// If the map is empty, it returns the zero value of K and false.
func (t *OrderedMap[K, V]) Min() (K, bool) {
	if t.IsEmpty() {
		var zero K
		return zero, false
	}
	return t.min(t.root).key, true
}

// Max returns the largest key in the OrderedMap and a boolean indicating success.
// If the map is empty, it returns the zero value of K and false.
func (t *OrderedMap[K, V]) Max() (K, bool) {
	if t.IsEmpty() {
		var zero K
		return zero, false
	}
	return t.max(t.root).key, true
}

// KeysInRange returns a slice of all keys in the OrderedMap between lo and hi, inclusive.
func (t *OrderedMap[K, V]) KeysInRange(lo, hi K) []K {
	queue := make([]K, 0)
	t.keysInRange(t.root, &queue, lo, hi)
	return queue
}

// isRed checks if a given node is red.
func (t *OrderedMap[K, V]) isRed(x *node[K, V]) bool {
	if x == nil {
		return false
	}
	return x.color == RED
}

// size returns the size of the subtree rooted at node x.
func (t *OrderedMap[K, V]) size(x *node[K, V]) int {
	if x == nil {
		return 0
	}
	return x.size
}

// get retrieves the value associated with the given key from the subtree rooted at x.
func (t *OrderedMap[K, V]) get(x *node[K, V], key K) (V, bool) {
	for x != nil {
		switch {
		case key < x.key:
			x = x.left
		case key > x.key:
			x = x.right
		default:
			return x.val, true
		}
	}
	var zero V
	return zero, false
}

// put inserts or updates a key-value pair in the subtree rooted at h.
func (t *OrderedMap[K, V]) put(h *node[K, V], key K, val V) *node[K, V] {
	if h == nil {
		return &node[K, V]{key: key, val: val, color: RED, size: 1}
	}

	switch {
	case key < h.key:
		h.left = t.put(h.left, key, val)
	case key > h.key:
		h.right = t.put(h.right, key, val)
	default:
		h.val = val
		return h
	}

	if t.isRed(h.right) && !t.isRed(h.left) {
		h = t.rotateLeft(h)
	}
	if t.isRed(h.left) && t.isRed(h.left.left) {
		h = t.rotateRight(h)
	}
	if t.isRed(h.left) && t.isRed(h.right) {
		t.flipColors(h)
	}

	h.size = t.size(h.left) + t.size(h.right) + 1
	return h
}

// rotateRight performs a right rotation on the given node.
func (t *OrderedMap[K, V]) rotateRight(h *node[K, V]) *node[K, V] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = t.size(h.left) + t.size(h.right) + 1
	return x
}

// rotateLeft performs a left rotation on the given node.
func (t *OrderedMap[K, V]) rotateLeft(h *node[K, V]) *node[K, V] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = t.size(h.left) + t.size(h.right) + 1
	return x
}

// flipColors flips the colors of a node and its two children.
func (t *OrderedMap[K, V]) flipColors(h *node[K, V]) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

// DeleteMin removes the smallest key and associated value from the map.
func (t *OrderedMap[K, V]) DeleteMin() {
	if t.IsEmpty() {
		panic("BST underflow")
	}

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}

	t.root = t.deleteMin(t.root)
	if !t.IsEmpty() {
		t.root.color = BLACK
	}
}

// deleteMin removes the node with the smallest key from the subtree rooted at h.
func (t *OrderedMap[K, V]) deleteMin(h *node[K, V]) *node[K, V] {
	if h.left == nil {
		return nil
	}

	if !t.isRed(h.left) && !t.isRed(h.left.left) {
		h = t.moveRedLeft(h)
	}

	h.left = t.deleteMin(h.left)
	return t.balance(h)
}

// DeleteMax removes the largest key and associated value from the map.
func (t *OrderedMap[K, V]) DeleteMax() {
	if t.IsEmpty() {
		panic("BST underflow")
	}

	if !t.isRed(t.root.left) && !t.isRed(t.root.right) {
		t.root.color = RED
	}

	t.root = t.deleteMax(t.root)
	if !t.IsEmpty() {
		t.root.color = BLACK
	}
}

// deleteMax removes the node with the largest key from the subtree rooted at h.
func (t *OrderedMap[K, V]) deleteMax(h *node[K, V]) *node[K, V] {
	if t.isRed(h.left) {
		h = t.rotateRight(h)
	}

	if h.right == nil {
		return nil
	}

	if !t.isRed(h.right) && !t.isRed(h.right.left) {
		h = t.moveRedRight(h)
	}

	h.right = t.deleteMax(h.right)

	return t.balance(h)
}

// delete removes the node with the given key from the subtree rooted at h.
func (t *OrderedMap[K, V]) delete(h *node[K, V], key K) *node[K, V] {
	if key < h.key {
		if !t.isRed(h.left) && !t.isRed(h.left.left) {
			h = t.moveRedLeft(h)
		}
		h.left = t.delete(h.left, key)
	} else {
		if t.isRed(h.left) {
			h = t.rotateRight(h)
		}
		if key == h.key && h.right == nil {
			return nil
		}
		if !t.isRed(h.right) && !t.isRed(h.right.left) {
			h = t.moveRedRight(h)
		}
		if key == h.key {
			x := t.min(h.right)
			h.key = x.key
			h.val = x.val
			h.right = t.deleteMin(h.right)
		} else {
			h.right = t.delete(h.right, key)
		}
	}
	return t.balance(h)
}

// moveRedLeft makes the left child or one of its children red.
func (t *OrderedMap[K, V]) moveRedLeft(h *node[K, V]) *node[K, V] {
	t.flipColors(h)
	if t.isRed(h.right.left) {
		h.right = t.rotateRight(h.right)
		h = t.rotateLeft(h)
		t.flipColors(h)
	}
	return h
}

// moveRedRight makes the right child or one of its children red.
func (t *OrderedMap[K, V]) moveRedRight(h *node[K, V]) *node[K, V] {
	t.flipColors(h)
	if t.isRed(h.left.left) {
		h = t.rotateRight(h)
		t.flipColors(h)
	}
	return h
}

// balance restores red-black tree invariants.
func (t *OrderedMap[K, V]) balance(h *node[K, V]) *node[K, V] {
	if t.isRed(h.right) && !t.isRed(h.left) {
		h = t.rotateLeft(h)
	}
	if t.isRed(h.left) && t.isRed(h.left.left) {
		h = t.rotateRight(h)
	}
	if t.isRed(h.left) && t.isRed(h.right) {
		t.flipColors(h)
	}

	h.size = t.size(h.left) + t.size(h.right) + 1
	return h
}

// min returns the node with the smallest key in the subtree rooted at x.
func (t *OrderedMap[K, V]) min(x *node[K, V]) *node[K, V] {
	if x.left == nil {
		return x
	}
	return t.min(x.left)
}

// max returns the node with the largest key in the subtree rooted at x.
func (t *OrderedMap[K, V]) max(x *node[K, V]) *node[K, V] {
	if x.right == nil {
		return x
	}
	return t.max(x.right)
}

// keysInRange collects keys in the given range [lo, hi] from the subtree rooted at x.
func (t *OrderedMap[K, V]) keysInRange(x *node[K, V], queue *[]K, lo, hi K) {
	if x == nil {
		return
	}
	// Claude had error in comparison order
	cmplt := lo < x.key
	cmple := lo <= x.key
	cmpgt := hi > x.key
	cmpge := hi >= x.key

	if cmplt {
		t.keysInRange(x.left, queue, lo, hi)
	}
	if cmple && cmpge {
		*queue = append(*queue, x.key)
	}
	if cmpgt {
		t.keysInRange(x.right, queue, lo, hi)
	}
}

func (t *OrderedMap[K, V]) keysInRangeBFS(x *node[K, V], queue *[]K) []K {

	if x == nil {
		return []K{}
	}

	// visit
	*queue = append(*queue, x.key)

	// go left
	t.keysInRangeBFS(x.left, queue)
	t.keysInRangeBFS(x.right, queue)

	return *queue
}

func (t *OrderedMap[K, V]) KeysInRangeBFS() []K {
	if t.IsEmpty() {
		return []K{}
	}

	queue := make([]K, 0)
	t.keysInRangeBFS(t.root, &queue)

	return queue
}
