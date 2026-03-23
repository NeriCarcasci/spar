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

    int countComponents(int n, std::vector<std::vector<int>>& edges) {
        parent.resize(n);
        rank_.resize(n, 0);
        std::iota(parent.begin(), parent.end(), 0);

        int components = n;
        for (auto& e : edges) {
            int px = find(e[0]), py = find(e[1]);
            if (px != py) {
                if (rank_[px] < rank_[py]) std::swap(px, py);
                parent[py] = px;
                if (rank_[px] == rank_[py]) rank_[px]++;
                components--;
            }
        }
        return components;
    }
};
