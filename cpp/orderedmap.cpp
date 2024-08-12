#include <iostream>
#include <stdexcept>
#include <vector>
#include <algorithm>
#include <memory>
#include <optional>
#include <functional>

template<typename Key, typename Value, typename Compare = std::less<Key>>
class OrderedMap {
private:
    static const bool RED = true;
    static const bool BLACK = false;

    struct Node {
        Key key;
        Value val;
        std::unique_ptr<Node> left;
        std::unique_ptr<Node> right;
        bool color;
        int size;

        Node(Key k, Value v, bool col, int sz) 
            : key(std::move(k)), val(std::move(v)), left(nullptr), right(nullptr), color(col), size(sz) {}
    };

    std::unique_ptr<Node> root;
    Compare comp;

    bool isRed(const Node* x) const noexcept {
        if (x == nullptr) return false;
        return x->color == RED;
    }

    int size(const Node* x) const noexcept {
        if (x == nullptr) return 0;
        return x->size;
    }

    std::optional<Value> get(const Node* x, const Key& key) const {
        while (x != nullptr) {
            if (comp(key, x->key)) x = x->left.get();
            else if (comp(x->key, key)) x = x->right.get();
            else return x->val;
        }
        return std::nullopt;
    }

    std::unique_ptr<Node> put(std::unique_ptr<Node> h, const Key& key, const Value& val) {
        if (h == nullptr) return std::make_unique<Node>(key, val, RED, 1);

        if (comp(key, h->key)) h->left = put(std::move(h->left), key, val);
        else if (comp(h->key, key)) h->right = put(std::move(h->right), key, val);
        else h->val = val;

        if (isRed(h->right.get()) && !isRed(h->left.get())) h = rotateLeft(std::move(h));
        if (isRed(h->left.get()) && isRed(h->left->left.get())) h = rotateRight(std::move(h));
        if (isRed(h->left.get()) && isRed(h->right.get())) flipColors(h.get());

        h->size = size(h->left.get()) + size(h->right.get()) + 1;
        return h;
    }

    std::unique_ptr<Node> rotateRight(std::unique_ptr<Node> h) {
        auto x = std::move(h->left);
        h->left = std::move(x->right);
        x->right = std::move(h);
        x->color = x->right->color;
        x->right->color = RED;
        x->size = x->right->size;
        x->right->size = size(x->right->left.get()) + size(x->right->right.get()) + 1;
        return x;
    }

    std::unique_ptr<Node> rotateLeft(std::unique_ptr<Node> h) {
        auto x = std::move(h->right);
        h->right = std::move(x->left);
        x->left = std::move(h);
        x->color = x->left->color;
        x->left->color = RED;
        x->size = x->left->size;
        x->left->size = size(x->left->left.get()) + size(x->left->right.get()) + 1;
        return x;
    }

    void flipColors(Node* h) noexcept {
        h->color = !h->color;
        h->left->color = !h->left->color;
        h->right->color = !h->right->color;
    }

    std::unique_ptr<Node> moveRedLeft(std::unique_ptr<Node> h) {
        flipColors(h.get());
        if (isRed(h->right->left.get())) {
            h->right = rotateRight(std::move(h->right));
            h = rotateLeft(std::move(h));
            flipColors(h.get());
        }
        return h;
    }

    std::unique_ptr<Node> moveRedRight(std::unique_ptr<Node> h) {
        flipColors(h.get());
        if (isRed(h->left->left.get())) {
            h = rotateRight(std::move(h));
            flipColors(h.get());
        }
        return h;
    }

    std::unique_ptr<Node> balance(std::unique_ptr<Node> h) {
        if (isRed(h->right.get()) && !isRed(h->left.get())) h = rotateLeft(std::move(h));
        if (isRed(h->left.get()) && isRed(h->left->left.get())) h = rotateRight(std::move(h));
        if (isRed(h->left.get()) && isRed(h->right.get())) flipColors(h.get());

        h->size = size(h->left.get()) + size(h->right.get()) + 1;
        return h;
    }

    std::unique_ptr<Node> deleteMin(std::unique_ptr<Node> h) {
        if (h->left == nullptr) return nullptr;

        if (!isRed(h->left.get()) && !isRed(h->left->left.get()))
            h = moveRedLeft(std::move(h));

        h->left = deleteMin(std::move(h->left));
        return balance(std::move(h));
    }

    std::unique_ptr<Node> deleteNode(std::unique_ptr<Node> h, const Key& key) {
        if (comp(key, h->key)) {
            if (!isRed(h->left.get()) && !isRed(h->left->left.get()))
                h = moveRedLeft(std::move(h));
            h->left = deleteNode(std::move(h->left), key);
        } else {
            if (isRed(h->left.get()))
                h = rotateRight(std::move(h));
            if (!comp(h->key, key) && !comp(key, h->key) && h->right == nullptr)
                return nullptr;
            if (!isRed(h->right.get()) && !isRed(h->right->left.get()))
                h = moveRedRight(std::move(h));
            if (!comp(h->key, key) && !comp(key, h->key)) {
                auto x = min(h->right.get());
                h->key = x->key;
                h->val = x->val;
                h->right = deleteMin(std::move(h->right));
            } else {
                h->right = deleteNode(std::move(h->right), key);
            }
        }
        return balance(std::move(h));
    }

    void keys(const Node* x, std::vector<Key>& v, const Key& lo, const Key& hi) const {
        if (x == nullptr) return;
        if (comp(lo, x->key)) keys(x->left.get(), v, lo, hi);
        if (!comp(x->key, lo) && !comp(hi, x->key)) v.push_back(x->key);
        if (comp(x->key, hi)) keys(x->right.get(), v, lo, hi);
    }

public:
    OrderedMap() : root(nullptr), comp(Compare()) {}
    explicit OrderedMap(const Compare& comp) : root(nullptr), comp(comp) {}

    OrderedMap(const OrderedMap& other) = delete;
    OrderedMap& operator=(const OrderedMap& other) = delete;

    OrderedMap(OrderedMap&& other) noexcept = default;
    OrderedMap& operator=(OrderedMap&& other) noexcept = default;

    int size() const noexcept {
        return size(root.get());
    }

    bool isEmpty() const noexcept {
        return root == nullptr;
    }

    std::optional<Value> get(const Key& key) const {
        return get(root.get(), key);
    }

    bool contains(const Key& key) const {
        return get(key).has_value();
    }

    void put(const Key& key, const Value& val) {
        root = put(std::move(root), key, val);
        root->color = BLACK;
    }

    void deleteMin() {
        if (isEmpty()) throw std::runtime_error("OrderedMap underflow");

        if (!isRed(root->left.get()) && !isRed(root->right.get()))
            root->color = RED;

        root = deleteMin(std::move(root));
        if (!isEmpty()) root->color = BLACK;
    }

    void deleteMax() {
        if (isEmpty()) throw std::runtime_error("OrderedMap underflow");

        if (!isRed(root->left.get()) && !isRed(root->right.get()))
            root->color = RED;

        root = deleteMax(std::move(root));
        if (!isEmpty()) root->color = BLACK;
    }

    void deleteNode(const Key& key) {
        if (!contains(key)) return;

        if (!isRed(root->left.get()) && !isRed(root->right.get()))
            root->color = RED;

        root = deleteNode(std::move(root), key);
        if (!isEmpty()) root->color = BLACK;
    }

    const Node* min(const Node* x) const {
        if (x->left == nullptr) return x;
        return min(x->left.get());
    }

    const Node* max(const Node* x) const {
        if (x->right == nullptr) return x;
        return max(x->right.get());
    }

    std::optional<Key> min() const {
        if (isEmpty()) return std::nullopt;
        return min(root.get())->key;
    }

    std::optional<Key> max() const {
        if (isEmpty()) return std::nullopt;
        return max(root.get())->key;
    }

    std::vector<Key> keys() const {
        if (isEmpty()) return std::vector<Key>();
        auto minKey = min();
        auto maxKey = max();
        if (minKey && maxKey) {
            return keys(*minKey, *maxKey);
        }
        return std::vector<Key>();
    }

    std::vector<Key> keys(const Key& lo, const Key& hi) const {
        std::vector<Key> v;
        keys(root.get(), v, lo, hi);
        return v;
    }
};

// Example usage
int main() {
    OrderedMap<std::string, int> map;

    map.put("C", 3);
    map.put("A", 1);
    map.put("G", 5);
    map.put("H", 6);
    map.put("B", 2);
    map.put("F", 4);

    std::cout << "Size: " << map.size() << std::endl;
    std::cout << "Contains 'B': " << (map.contains("B") ? "Yes" : "No") << std::endl;
    
    if (auto value = map.get("C"); value) {
        std::cout << "Value of 'C': " << *value << std::endl;
    }

    map.deleteNode("B");

    std::cout << "Size after deleting 'B': " << map.size() << std::endl;

    std::cout << "Keys: \n";
    for (const auto& key : map.keys()) {
        if (auto value = map.get(key); value) {
            std::cout << key << " " << *value << "\n";
        }
    }
    std::cout << std::endl;

    return 0;
}
