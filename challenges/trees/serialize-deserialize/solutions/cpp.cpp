#include <string>
#include <sstream>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

class Solution {
public:
    std::string serialize(TreeNode* root) {
        if (!root) return "N";
        return std::to_string(root->val) + "," + serialize(root->left) + "," + serialize(root->right);
    }

    TreeNode* deserializeHelper(std::istringstream& ss) {
        std::string token;
        std::getline(ss, token, ',');
        if (token == "N") return nullptr;
        TreeNode* node = new TreeNode(std::stoi(token));
        node->left = deserializeHelper(ss);
        node->right = deserializeHelper(ss);
        return node;
    }

    TreeNode* deserialize(std::string data) {
        std::istringstream ss(data);
        return deserializeHelper(ss);
    }
};
