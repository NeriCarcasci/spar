#include <vector>
#include <numeric>

class Solution {
public:
    std::vector<int> parent;
    std::vector<int> rank_;

    int find(int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    }

    bool unite(int x, int y) {
        int px = find(x), py = find(y);
        if (px == py) return false;
        if (rank_[px] < rank_[py]) std::swap(px, py);
        parent[py] = px;
        if (rank_[px] == rank_[py]) rank_[px]++;
        return true;
    }

    bool validTree(int n, std::vector<std::vector<int>>& edges) {
        if (static_cast<int>(edges.size()) != n - 1) return false;
        parent.resize(n);
        rank_.resize(n, 0);
        std::iota(parent.begin(), parent.end(), 0);

        for (auto& e : edges) {
            if (!unite(e[0], e[1])) return false;
        }
        return true;
    }
};
