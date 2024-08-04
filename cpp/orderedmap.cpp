#include <iostream>
#include <stdexcept>
#include <queue>
#include <algorithm>

template<typename Key, typename Value>
class OrderedMap {
private:
    static const bool RED = true;
    static const bool BLACK = false;

    struct Node {
        Key key;
        Value val;
        Node* left;
        Node* right;
        bool color;
        int size;

        Node(Key k, Value v, bool col, int sz) 
            : key(k), val(v), left(nullptr), right(nullptr), color(col), size(sz) {}
    };

    Node* root;

    bool isRed(Node* x) const {
        if (x == nullptr) return false;
        return x->color == RED;
    }

    int size(Node* x) const {
        if (x == nullptr) return 0;
        return x->size;
    }

    Value get(Node* x, const Key& key) const {
        while (x != nullptr) {
            if (key < x->key) x = x->left;
            else if (x->key < key) x = x->right;
            else return x->val;
        }
        throw std::runtime_error("Key not found");
    }

    Node* put(Node* h, const Key& key, const Value& val) {
        if (h == nullptr) return new Node(key, val, RED, 1);

        if (key < h->key) h->left = put(h->left, key, val);
        else if (h->key < key) h->right = put(h->right, key, val);
        else h->val = val;

        if (isRed(h->right) && !isRed(h->left)) h = rotateLeft(h);
        if (isRed(h->left) && isRed(h->left->left)) h = rotateRight(h);
        if (isRed(h->left) && isRed(h->right)) flipColors(h);

        h->size = size(h->left) + size(h->right) + 1;
        return h;
    }

    Node* rotateRight(Node* h) {
        Node* x = h->left;
        h->left = x->right;
        x->right = h;
        x->color = h->color;
        h->color = RED;
        x->size = h->size;
        h->size = size(h->left) + size(h->right) + 1;
        return x;
    }

    Node* rotateLeft(Node* h) {
        Node* x = h->right;
        h->right = x->left;
        x->left = h;
        x->color = h->color;
        h->color = RED;
        x->size = h->size;
        h->size = size(h->left) + size(h->right) + 1;
        return x;
    }

    void flipColors(Node* h) {
        h->color = !h->color;
        h->left->color = !h->left->color;
        h->right->color = !h->right->color;
    }

    Node* moveRedLeft(Node* h) {
        flipColors(h);
        if (isRed(h->right->left)) {
            h->right = rotateRight(h->right);
            h = rotateLeft(h);
            flipColors(h);
        }
        return h;
    }

    Node* moveRedRight(Node* h) {
        flipColors(h);
        if (isRed(h->left->left)) {
            h = rotateRight(h);
            flipColors(h);
        }
        return h;
    }

    Node* balance(Node* h) {
        if (isRed(h->right) && !isRed(h->left)) h = rotateLeft(h);
        if (isRed(h->left) && isRed(h->left->left)) h = rotateRight(h);
        if (isRed(h->left) && isRed(h->right)) flipColors(h);

        h->size = size(h->left) + size(h->right) + 1;
        return h;
    }

    Node* min(Node* x) const {
        if (x->left == nullptr) return x;
        return min(x->left);
    }

    Node* deleteMin(Node* h) {
        if (h->left == nullptr) return nullptr;

        if (!isRed(h->left) && !isRed(h->left->left))
            h = moveRedLeft(h);

        h->left = deleteMin(h->left);
        return balance(h);
    }

    Node* deleteNode(Node* h, const Key& key) {
        if (key < h->key) {
            if (!isRed(h->left) && !isRed(h->left->left))
                h = moveRedLeft(h);
            h->left = deleteNode(h->left, key);
        } else {
            if (isRed(h->left))
                h = rotateRight(h);
            if (key == h->key && h->right == nullptr)
                return nullptr;
            if (!isRed(h->right) && !isRed(h->right->left))
                h = moveRedRight(h);
            if (key == h->key) {
                Node* x = min(h->right);
                h->key = x->key;
                h->val = x->val;
                h->right = deleteMin(h->right);
            } else {
                h->right = deleteNode(h->right, key);
            }
        }
        return balance(h);
    }

    void keys(Node* x, std::queue<Key>& queue, const Key& lo, const Key& hi) const {
        if (x == nullptr) return;
        if (lo < x->key) keys(x->left, queue, lo, hi);
        if (lo <= x->key && x->key <= hi) queue.push(x->key);
        if (x->key < hi) keys(x->right, queue, lo, hi);
    }

public:
    OrderedMap() : root(nullptr) {}

    int size() const {
        return size(root);
    }

    bool isEmpty() const {
        return root == nullptr;
    }

    Value get(const Key& key) const {
        return get(root, key);
    }

    bool contains(const Key& key) const {
        try {
            get(key);
            return true;
        } catch (const std::runtime_error&) {
            return false;
        }
    }

    void put(const Key& key, const Value& val) {
        root = put(root, key, val);
        root->color = BLACK;
    }

    void deleteMin() {
        if (isEmpty()) throw std::runtime_error("OrderedMap underflow");

        if (!isRed(root->left) && !isRed(root->right))
            root->color = RED;

        root = deleteMin(root);
        if (!isEmpty()) root->color = BLACK;
    }

    void deleteMax() {
        if (isEmpty()) throw std::runtime_error("OrderedMap underflow");

        if (!isRed(root->left) && !isRed(root->right))
            root->color = RED;

        root = deleteMax(root);
        if (!isEmpty()) root->color = BLACK;
    }

    void deleteNode(const Key& key) {
        if (!contains(key)) return;

        if (!isRed(root->left) && !isRed(root->right))
            root->color = RED;

        root = deleteNode(root, key);
        if (!isEmpty()) root->color = BLACK;
    }

    Key min() const {
        if (isEmpty()) throw std::runtime_error("Called min() with empty OrderedMap");
        return min(root)->key;
    }

    Key max() const {
        if (isEmpty()) throw std::runtime_error("Called max() with empty OrderedMap");
        return max(root)->key;
    }

    std::vector<Key> keys() const {
        if (isEmpty()) return std::vector<Key>();
        return keys(min(), max());
    }

    std::vector<Key> keys(const Key& lo, const Key& hi) const {
        std::queue<Key> queue;
        keys(root, queue, lo, hi);
        return std::vector<Key>(queue.front(), queue.back());
    }
};

// Example usage
int main() {
    OrderedMap<std::string, int> map;

    map.put("A", 1);
    map.put("B", 2);
    map.put("C", 3);

    std::cout << "Size: " << map.size() << std::endl;
    std::cout << "Contains 'B': " << (map.contains("B") ? "Yes" : "No") << std::endl;
    std::cout << "Value of 'C': " << map.get("C") << std::endl;

    map.deleteNode("B");

    std::cout << "Size after deleting 'B': " << map.size() << std::endl;

    std::cout << "Keys: ";
    for (const auto& key : map.keys()) {
        std::cout << key << " ";
    }
    std::cout << std::endl;

    return 0;
}
