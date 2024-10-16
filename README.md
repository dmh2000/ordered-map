# ordered-map

Go has a built in 'map' datastructure that provides a key/value store. It is implemented with an unordered hashmap. A hashmap is good for many (most?) usage of a key/value store because its runtime complexity is O(1). The drawback of the go 'map' is that it is not ordered. Its can't be iterated over in order by keys. One solution is to implement an ordered map using a Red-Black search tree (or an AVL tree). These trees are O(log N) search, insert and delete with O(N) memory. They do this by a resonably tricky algorithm that keeps the tree mostly balanced at all times. Both types usupport iterating in order.

The reference for the algorithm is [Algorithms, 4th edition by Robert Sedgewick and Kevin Wayne, Addison-Wesley Professional, 2011, ISBN 0-321-57351-X](http://algs4.cs.princeton.edu). This site has code for many algorithms and data structures, along with a textbook and a video course.  The Go implementation here is based on [the original Java code from ](https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/RedBlackBST.java). 

The original port was made automatically using [Anthropic Claude-3.5-Sonnet](https://www.anthropic.com/) and [aider-chat](https://aider.chat/). The Java code is very extensive and complete, and the Go version is a manually pared down subset that supports functions similar to what a Go map provides, including Get, Put, Delete, Contains, IsEmpty and iterate over (in order). 

## Rust Implementation

A Rust implementation of the `OrderedMap` is also available. The Rust version provides similar functionality to the Go version, including methods for getting, putting, deleting, and checking the existence of keys, as well as iterating over keys in order.

### Usage

To use the Rust implementation, add the following to your `Cargo.toml`:

```toml
[dependencies]
orderedmap = { path = "path/to/orderedmap" }
```

Then, in your Rust code, you can use the `OrderedMap` as follows:

```rust
use orderedmap::OrderedMap;

fn main() {
    let mut map = OrderedMap::new();

    map.put("C", 3);
    map.put("A", 1);
    map.put("G", 5);
    map.put("H", 6);
    map.put("B", 2);
    map.put("F", 4);

    println!("Size: {}", map.size());
    println!("Contains 'B': {}", map.contains("B"));
    println!("Value of 'C': {:?}", map.get("C"));

    map.delete("B");

    println!("Size after deleting 'B': {}", map.size());

    println!("Keys: ");
    for key in map.keys() {
        println!("{}: {:?}", key, map.get(key));
    }
}
```
