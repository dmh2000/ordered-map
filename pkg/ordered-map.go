package pkg

import (
	"fmt"
)

type color bool

const (
	RED   color = true
	BLACK color = false
)

type node[K comparable, V any] struct {
	key         K
	val         V
	left, right *node[K, V]
	color       color
	size        int
}

type RedBlackBST[K comparable, V any] struct {
	root *node[K, V]
}

func NewRedBlackBST[K comparable, V any]() *RedBlackBST[K, V] {
	return &RedBlackBST[K, V]{}
}

func (t *RedBlackBST[K, V]) isRed(x *node[K, V]) bool {
	if x == nil {
		return false
	}
	return x.color == RED
}

func (t *RedBlackBST[K, V]) size(x *node[K, V]) int {
	if x == nil {
		return 0
	}
	return x.size
}

func (t *RedBlackBST[K, V]) Size() int {
	return t.size(t.root)
}

func (t *RedBlackBST[K, V]) IsEmpty() bool {
	return t.root == nil
}

func (t *RedBlackBST[K, V]) Get(key K) (V, bool) {
	return t.get(t.root, key)
}

func (t *RedBlackBST[K, V]) get(x *node[K, V], key K) (V, bool) {
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

func (t *RedBlackBST[K, V]) Contains(key K) bool {
	_, found := t.Get(key)
	return found
}

func (t *RedBlackBST[K, V]) Put(key K, val V) {
	t.root = t.put(t.root, key, val)
	t.root.color = BLACK
}

func (t *RedBlackBST[K, V]) put(h *node[K, V], key K, val V) *node[K, V] {
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

func (t *RedBlackBST[K, V]) rotateRight(h *node[K, V]) *node[K, V] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = t.size(h.left) + t.size(h.right) + 1
	return x
}

func (t *RedBlackBST[K, V]) rotateLeft(h *node[K, V]) *node[K, V] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = RED
	x.size = h.size
	h.size = t.size(h.left) + t.size(h.right) + 1
	return x
}

func (t *RedBlackBST[K, V]) flipColors(h *node[K, V]) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func (t *RedBlackBST[K, V]) DeleteMin() {
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

func (t *RedBlackBST[K, V]) deleteMin(h *node[K, V]) *node[K, V] {
	if h.left == nil {
		return nil
	}

	if !t.isRed(h.left) && !t.isRed(h.left.left) {
		h = t.moveRedLeft(h)
	}

	h.left = t.deleteMin(h.left)
	return t.balance(h)
}

func (t *RedBlackBST[K, V]) DeleteMax() {
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

func (t *RedBlackBST[K, V]) deleteMax(h *node[K, V]) *node[K, V] {
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

func (t *RedBlackBST[K, V]) Delete(key K) {
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

func (t *RedBlackBST[K, V]) delete(h *node[K, V], key K) *node[K, V] {
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

func (t *RedBlackBST[K, V]) moveRedLeft(h *node[K, V]) *node[K, V] {
	t.flipColors(h)
	if t.isRed(h.right.left) {
		h.right = t.rotateRight(h.right)
		h = t.rotateLeft(h)
		t.flipColors(h)
	}
	return h
}

func (t *RedBlackBST[K, V]) moveRedRight(h *node[K, V]) *node[K, V] {
	t.flipColors(h)
	if t.isRed(h.left.left) {
		h = t.rotateRight(h)
		t.flipColors(h)
	}
	return h
}

func (t *RedBlackBST[K, V]) balance(h *node[K, V]) *node[K, V] {
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

func (t *RedBlackBST[K, V]) Height() int {
	return t.height(t.root)
}

func (t *RedBlackBST[K, V]) height(x *node[K, V]) int {
	if x == nil {
		return -1
	}
	return 1 + max(t.height(x.left), t.height(x.right))
}

func (t *RedBlackBST[K, V]) Min() (K, bool) {
	if t.IsEmpty() {
		var zero K
		return zero, false
	}
	return t.min(t.root).key, true
}

func (t *RedBlackBST[K, V]) min(x *node[K, V]) *node[K, V] {
	if x.left == nil {
		return x
	}
	return t.min(x.left)
}

func (t *RedBlackBST[K, V]) Max() (K, bool) {
	if t.IsEmpty() {
		var zero K
		return zero, false
	}
	return t.max(t.root).key, true
}

func (t *RedBlackBST[K, V]) max(x *node[K, V]) *node[K, V] {
	if x.right == nil {
		return x
	}
	return t.max(x.right)
}

func (t *RedBlackBST[K, V]) Floor(key K) (K, bool) {
	x := t.floor(t.root, key)
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

func (t *RedBlackBST[K, V]) floor(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}
	if key == x.key {
		return x
	}
	if key < x.key {
		return t.floor(x.left, key)
	}
	y := t.floor(x.right, key)
	if y != nil {
		return y
	}
	return x
}

func (t *RedBlackBST[K, V]) Ceiling(key K) (K, bool) {
	x := t.ceiling(t.root, key)
	if x == nil {
		var zero K
		return zero, false
	}
	return x.key, true
}

func (t *RedBlackBST[K, V]) ceiling(x *node[K, V], key K) *node[K, V] {
	if x == nil {
		return nil
	}
	if key == x.key {
		return x
	}
	if key > x.key {
		return t.ceiling(x.right, key)
	}
	y := t.ceiling(x.left, key)
	if y != nil {
		return y
	}
	return x
}

func (t *RedBlackBST[K, V]) Select(rank int) (K, bool) {
	if rank < 0 || rank >= t.Size() {
		var zero K
		return zero, false
	}
	return t.select(t.root, rank).key, true
}

func (t *RedBlackBST[K, V]) select(x *node[K, V], rank int) *node[K, V] {
	if x == nil {
		return nil
	}
	leftSize := t.size(x.left)
	if leftSize > rank {
		return t.select(x.left, rank)
	} else if leftSize < rank {
		return t.select(x.right, rank-leftSize-1)
	} else {
		return x
	}
}

func (t *RedBlackBST[K, V]) Rank(key K) int {
	return t.rank(key, t.root)
}

func (t *RedBlackBST[K, V]) rank(key K, x *node[K, V]) int {
	if x == nil {
		return 0
	}
	if key < x.key {
		return t.rank(key, x.left)
	} else if key > x.key {
		return 1 + t.size(x.left) + t.rank(key, x.right)
	} else {
		return t.size(x.left)
	}
}

func (t *RedBlackBST[K, V]) Keys() []K {
	if t.IsEmpty() {
		return []K{}
	}
	min, _ := t.Min()
	max, _ := t.Max()
	return t.KeysInRange(min, max)
}

func (t *RedBlackBST[K, V]) KeysInRange(lo, hi K) []K {
	queue := make([]K, 0)
	t.keysInRange(t.root, &queue, lo, hi)
	return queue
}

func (t *RedBlackBST[K, V]) keysInRange(x *node[K, V], queue *[]K, lo, hi K) {
	if x == nil {
		return
	}
	cmplo := lo < x.key
	cmphi := hi > x.key
	if cmplo {
		t.keysInRange(x.left, queue, lo, hi)
	}
	if cmplo && cmphi {
		*queue = append(*queue, x.key)
	}
	if cmphi {
		t.keysInRange(x.right, queue, lo, hi)
	}
}

func (t *RedBlackBST[K, V]) SizeInRange(lo, hi K) int {
	if lo > hi {
		return 0
	}
	if t.Contains(hi) {
		return t.Rank(hi) - t.Rank(lo) + 1
	} else {
		return t.Rank(hi) - t.Rank(lo)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
