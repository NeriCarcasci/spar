#include <algorithm>
#include <climits>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

class Solution {
public:
    int result;

    int dfs(TreeNode* node) {
        if (!node) return 0;
        int leftGain = std::max(dfs(node->left), 0);
        int rightGain = std::max(dfs(node->right), 0);
        result = std::max(result, node->val + leftGain + rightGain);
        return node->val + std::max(leftGain, rightGain);
    }

    int maxPathSum(TreeNode* root) {
        result = INT_MIN;
        dfs(root);
        return result;
    }
};
