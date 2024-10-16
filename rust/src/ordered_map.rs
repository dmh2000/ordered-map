use std::cmp::Ordering;

#[derive(Debug)]
enum Color {
    Red,
    Black,
}

#[derive(Debug)]
struct Node<K, V> {
    key: K,
    value: V,
    left: Option<Box<Node<K, V>>>,
    right: Option<Box<Node<K, V>>>,
    color: Color,
    size: usize,
}

impl<K, V> Node<K, V> {
    fn new(key: K, value: V, color: Color, size: usize) -> Self {
        Node {
            key,
            value,
            left: None,
            right: None,
            color,
            size,
        }
    }
}

pub struct OrderedMap<K, V> {
    root: Option<Box<Node<K, V>>>,
}

impl<K: Ord, V> OrderedMap<K, V> {
    pub fn new() -> Self {
        OrderedMap { root: None }
    }

    pub fn is_empty(&self) -> bool {
        self.size() == 0
    }

    pub fn size(&self) -> usize {
        self.size_node(&self.root)
    }

    fn size_node(&self, node: &Option<Box<Node<K, V>>>) -> usize {
        match node {
            Some(n) => n.size,
            None => 0,
        }
    }

    fn is_red(node: &Option<Box<Node<K, V>>>) -> bool {
        match node {
            Some(n) => matches!(n.color, Color::Red),
            None => false,
        }
    }

    pub fn get(&self, key: K) -> Option<&V> {
        self.get_node(&self.root, key)
    }

    fn get_node<'a>(&'a self, node: &'a Option<Box<Node<K, V>>>, key: K) -> Option<&'a V> {
        match node {
            Some(n) => match key.cmp(&n.key) {
                Ordering::Less => self.get_node(&n.left, key),
                Ordering::Greater => self.get_node(&n.right, key),
                Ordering::Equal => Some(&n.value),
            },
            None => None,
        }
    }

    pub fn put(&mut self, key: K, value: V) {
        self.root = Self::put_node(self.root.take(), key, value);
        if let Some(ref mut root) = self.root {
            root.color = Color::Black;
        }
    }

    fn put_node(node: Option<Box<Node<K, V>>>, key: K, value: V) -> Option<Box<Node<K, V>>> {
        let mut h = match node {
            Some(n) => n,
            None => return Some(Box::new(Node::new(key, value, Color::Red, 1))),
        };

        match key.cmp(&h.key) {
            Ordering::Less => h.left = Self::put_node(h.left.take(), key, value),
            Ordering::Greater => h.right = Self::put_node(h.right.take(), key, value),
            Ordering::Equal => h.value = value,
        }

        if Self::is_red(&h.right) && !Self::is_red(&h.left) {
            h = Self::rotate_left(h);
        }
        if Self::is_red(&h.left) && Self::is_red(&h.left.as_ref().unwrap().left) {
            h = Self::rotate_right(h);
        }
        if Self::is_red(&h.left) && Self::is_red(&h.right) {
            Self::flip_colors(&mut h);
        }

        h.size = Self::size_node(&h.left) + Self::size_node(&h.right) + 1;
        Some(h)
    }

    fn rotate_left(mut h: Box<Node<K, V>>) -> Box<Node<K, V>> {
        let mut x = h.right.take().unwrap();
        h.right = x.left.take();
        x.left = Some(h);
        x.color = x.left.as_ref().unwrap().color.clone();
        x.left.as_mut().unwrap().color = Color::Red;
        x.size = x.left.as_ref().unwrap().size;
        x.left.as_mut().unwrap().size = Self::size_node(&x.left.as_ref().unwrap().left)
            + Self::size_node(&x.left.as_ref().unwrap().right)
            + 1;
        x
    }

    fn rotate_right(mut h: Box<Node<K, V>>) -> Box<Node<K, V>> {
        let mut x = h.left.take().unwrap();
        h.left = x.right.take();
        x.right = Some(h);
        x.color = x.right.as_ref().unwrap().color.clone();
        x.right.as_mut().unwrap().color = Color::Red;
        x.size = x.right.as_ref().unwrap().size;
        x.right.as_mut().unwrap().size = Self::size_node(&x.right.as_ref().unwrap().left)
            + Self::size_node(&x.right.as_ref().unwrap().right)
            + 1;
        x
    }

    fn flip_colors(h: &mut Box<Node<K, V>>) {
        h.color = match h.color {
            Color::Red => Color::Black,
            Color::Black => Color::Red,
        };
        if let Some(ref mut left) = h.left {
            left.color = match left.color {
                Color::Red => Color::Black,
                Color::Black => Color::Red,
            };
        }
        if let Some(ref mut right) = h.right {
            right.color = match right.color {
                Color::Red => Color::Black,
                Color::Black => Color::Red,
            };
        }
    }

    pub fn delete_min(&mut self) {
        if self.is_empty() {
            panic!("OrderedMap underflow");
        }

        if !Self::is_red(&self.root.as_ref().unwrap().left)
            && !Self::is_red(&self.root.as_ref().unwrap().right)
        {
            self.root.as_mut().unwrap().color = Color::Red;
        }

        self.root = Self::delete_min_node(self.root.take());
        if !self.is_empty() {
            self.root.as_mut().unwrap().color = Color::Black;
        }
    }

    fn delete_min_node(mut h: Option<Box<Node<K, V>>>) -> Option<Box<Node<K, V>>> {
        if h.as_ref().unwrap().left.is_none() {
            return None;
        }

        if !Self::is_red(&h.as_ref().unwrap().left)
            && !Self::is_red(&h.as_ref().unwrap().left.as_ref().unwrap().left)
        {
            h = Some(Self::move_red_left(h.unwrap()));
        }

        h.as_mut().unwrap().left = Self::delete_min_node(h.as_mut().unwrap().left.take());
        Some(Self::balance(h.unwrap()))
    }

    pub fn delete_max(&mut self) {
        if self.is_empty() {
            panic!("OrderedMap underflow");
        }

        if !Self::is_red(&self.root.as_ref().unwrap().left)
            && !Self::is_red(&self.root.as_ref().unwrap().right)
        {
            self.root.as_mut().unwrap().color = Color::Red;
        }

        self.root = Self::delete_max_node(self.root.take());
        if !self.is_empty() {
            self.root.as_mut().unwrap().color = Color::Black;
        }
    }

    fn delete_max_node(mut h: Option<Box<Node<K, V>>>) -> Option<Box<Node<K, V>>> {
        if Self::is_red(&h.as_ref().unwrap().left) {
            h = Some(Self::rotate_right(h.unwrap()));
        }

        if h.as_ref().unwrap().right.is_none() {
            return None;
        }

        if !Self::is_red(&h.as_ref().unwrap().right)
            && !Self::is_red(&h.as_ref().unwrap().right.as_ref().unwrap().left)
        {
            h = Some(Self::move_red_right(h.unwrap()));
        }

        h.as_mut().unwrap().right = Self::delete_max_node(h.as_mut().unwrap().right.take());
        Some(Self::balance(h.unwrap()))
    }

    pub fn delete(&mut self, key: K) {
        if !self.contains(&key) {
            return;
        }

        if !Self::is_red(&self.root.as_ref().unwrap().left)
            && !Self::is_red(&self.root.as_ref().unwrap().right)
        {
            self.root.as_mut().unwrap().color = Color::Red;
        }

        self.root = Self::delete_node(self.root.take(), key);
        if !self.is_empty() {
            self.root.as_mut().unwrap().color = Color::Black;
        }
    }

    fn delete_node(mut h: Option<Box<Node<K, V>>>, key: K) -> Option<Box<Node<K, V>>> {
        if key < h.as_ref().unwrap().key {
            if !Self::is_red(&h.as_ref().unwrap().left)
                && !Self::is_red(&h.as_ref().unwrap().left.as_ref().unwrap().left)
            {
                h = Some(Self::move_red_left(h.unwrap()));
            }
            h.as_mut().unwrap().left = Self::delete_node(h.as_mut().unwrap().left.take(), key);
        } else {
            if Self::is_red(&h.as_ref().unwrap().left) {
                h = Some(Self::rotate_right(h.unwrap()));
            }
            if key == h.as_ref().unwrap().key && h.as_ref().unwrap().right.is_none() {
                return None;
            }
            if !Self::is_red(&h.as_ref().unwrap().right)
                && !Self::is_red(&h.as_ref().unwrap().right.as_ref().unwrap().left)
            {
                h = Some(Self::move_red_right(h.unwrap()));
            }
            if key == h.as_ref().unwrap().key {
                let x = Self::min_node(h.as_ref().unwrap().right.as_ref().unwrap());
                h.as_mut().unwrap().key = x.key.clone();
                h.as_mut().unwrap().value = x.value.clone();
                h.as_mut().unwrap().right =
                    Self::delete_min_node(h.as_mut().unwrap().right.take());
            } else {
                h.as_mut().unwrap().right = Self::delete_node(h.as_mut().unwrap().right.take(), key);
            }
        }
        Some(Self::balance(h.unwrap()))
    }

    fn move_red_left(mut h: Box<Node<K, V>>) -> Box<Node<K, V>> {
        Self::flip_colors(&mut h);
        if Self::is_red(&h.right.as_ref().unwrap().left) {
            h.right = Some(Self::rotate_right(h.right.take().unwrap()));
            h = Self::rotate_left(h);
            Self::flip_colors(&mut h);
        }
        h
    }

    fn move_red_right(mut h: Box<Node<K, V>>) -> Box<Node<K, V>> {
        Self::flip_colors(&mut h);
        if Self::is_red(&h.left.as_ref().unwrap().left) {
            h = Self::rotate_right(h);
            Self::flip_colors(&mut h);
        }
        h
    }

    fn balance(mut h: Box<Node<K, V>>) -> Box<Node<K, V>> {
        if Self::is_red(&h.right) && !Self::is_red(&h.left) {
            h = Self::rotate_left(h);
        }
        if Self::is_red(&h.left) && Self::is_red(&h.left.as_ref().unwrap().left) {
            h = Self::rotate_right(h);
        }
        if Self::is_red(&h.left) && Self::is_red(&h.right) {
            Self::flip_colors(&mut h);
        }

        h.size = Self::size_node(&h.left) + Self::size_node(&h.right) + 1;
        h
    }

    fn min_node(x: &Box<Node<K, V>>) -> &Box<Node<K, V>> {
        match &x.left {
            Some(left) => Self::min_node(left),
            None => x,
        }
    }

    fn max_node(x: &Box<Node<K, V>>) -> &Box<Node<K, V>> {
        match &x.right {
            Some(right) => Self::max_node(right),
            None => x,
        }
    }

    pub fn min(&self) -> Option<&K> {
        self.root.as_ref().map(|node| &Self::min_node(node).key)
    }

    pub fn max(&self) -> Option<&K> {
        self.root.as_ref().map(|node| &Self::max_node(node).key)
    }

    pub fn keys(&self) -> Vec<&K> {
        let mut vec = Vec::new();
        self.keys_in_order(&self.root, &mut vec);
        vec
    }

    fn keys_in_order<'a>(&'a self, node: &'a Option<Box<Node<K, V>>>, vec: &mut Vec<&'a K>) {
        if let Some(n) = node {
            self.keys_in_order(&n.left, vec);
            vec.push(&n.key);
            self.keys_in_order(&n.right, vec);
        }
    }
}
