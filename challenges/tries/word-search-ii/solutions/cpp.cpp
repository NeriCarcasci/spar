#include <vector>
#include <string>
#include <algorithm>

struct TrieNode {
    TrieNode* children[26] = {};
    std::string word;
};

class Solution {
public:
    std::vector<std::string> findWords(std::vector<std::vector<char>>& board, std::vector<std::string>& words) {
        TrieNode* root = new TrieNode();
        for (auto& w : words) {
            TrieNode* node = root;
            for (char ch : w) {
                int idx = ch - 'a';
                if (!node->children[idx]) node->children[idx] = new TrieNode();
                node = node->children[idx];
            }
            node->word = w;
        }

        int rows = board.size(), cols = board[0].size();
        std::vector<std::string> result;
        int dirs[4][2] = {{0,1},{0,-1},{1,0},{-1,0}};

        std::function<void(int, int, TrieNode*)> dfs = [&](int r, int c, TrieNode* node) {
            char ch = board[r][c];
            if (ch == '#') return;
            TrieNode* next = node->children[ch - 'a'];
            if (!next) return;
            if (!next->word.empty()) {
                result.push_back(next->word);
                next->word.clear();
            }
            board[r][c] = '#';
            for (auto& d : dirs) {
                int nr = r + d[0], nc = c + d[1];
                if (nr >= 0 && nr < rows && nc >= 0 && nc < cols) {
                    dfs(nr, nc, next);
                }
            }
            board[r][c] = ch;
        };

        for (int r = 0; r < rows; r++) {
            for (int c = 0; c < cols; c++) {
                dfs(r, c, root);
            }
        }
        std::sort(result.begin(), result.end());
        return result;
    }
};
