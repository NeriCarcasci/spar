#include <vector>
#include <unordered_map>

class Node {
public:
    int val;
    std::vector<Node*> neighbors;
    Node() : val(0), neighbors(std::vector<Node*>()) {}
    Node(int _val) : val(_val), neighbors(std::vector<Node*>()) {}
    Node(int _val, std::vector<Node*> _neighbors) : val(_val), neighbors(_neighbors) {}
};

class Solution {
public:
    std::unordered_map<Node*, Node*> cloned;

    Node* cloneGraph(Node* node) {
        if (!node) return nullptr;
        if (cloned.count(node)) return cloned[node];
        Node* copy = new Node(node->val);
        cloned[node] = copy;
        for (Node* neighbor : node->neighbors) {
            copy->neighbors.push_back(cloneGraph(neighbor));
        }
        return copy;
    }
};
