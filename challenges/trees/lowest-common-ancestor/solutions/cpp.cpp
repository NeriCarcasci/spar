struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

class Solution {
public:
    int lowestCommonAncestor(TreeNode* root, int p, int q) {
        TreeNode* cur = root;
        while (cur) {
            if (p < cur->val && q < cur->val) {
                cur = cur->left;
            } else if (p > cur->val && q > cur->val) {
                cur = cur->right;
            } else {
                return cur->val;
            }
        }
        return -1;
    }
};
