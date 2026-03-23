#include <string>
#include <array>

struct TrieNode {
    std::array<TrieNode*, 26> children{};
    bool isEnd = false;
};

class Trie {
    TrieNode* root;
public:
    Trie() : root(new TrieNode()) {}

    void insert(const std::string& word) {
        TrieNode* node = root;
        for (char ch : word) {
            int idx = ch - 'a';
            if (!node->children[idx]) {
                node->children[idx] = new TrieNode();
            }
            node = node->children[idx];
        }
        node->isEnd = true;
    }

    bool search(const std::string& word) {
        TrieNode* node = find(word);
        return node && node->isEnd;
    }

    bool startsWith(const std::string& prefix) {
        return find(prefix) != nullptr;
    }

private:
    TrieNode* find(const std::string& prefix) {
        TrieNode* node = root;
        for (char ch : prefix) {
            int idx = ch - 'a';
            if (!node->children[idx]) return nullptr;
            node = node->children[idx];
        }
        return node;
    }
};
