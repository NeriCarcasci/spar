#include <string>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

class Solution {
public:
    std::string serialize(TreeNode* root) {
        return "";
    }

    TreeNode* deserialize(std::string data) {
        return nullptr;
    }
};
