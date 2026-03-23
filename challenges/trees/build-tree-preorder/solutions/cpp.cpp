#include <vector>
#include <unordered_map>

struct TreeNode {
    int val;
    TreeNode* left;
    TreeNode* right;
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
};

class Solution {
public:
    std::unordered_map<int, int> inMap;
    int preIdx;

    TreeNode* build(std::vector<int>& preorder, int inLeft, int inRight) {
        if (inLeft > inRight) return nullptr;
        int rootVal = preorder[preIdx++];
        TreeNode* root = new TreeNode(rootVal);
        int mid = inMap[rootVal];
        root->left = build(preorder, inLeft, mid - 1);
        root->right = build(preorder, mid + 1, inRight);
        return root;
    }

    TreeNode* buildTree(std::vector<int>& preorder, std::vector<int>& inorder) {
        preIdx = 0;
        inMap.clear();
        for (int i = 0; i < static_cast<int>(inorder.size()); i++) {
            inMap[inorder[i]] = i;
        }
        return build(preorder, 0, static_cast<int>(inorder.size()) - 1);
    }
};
