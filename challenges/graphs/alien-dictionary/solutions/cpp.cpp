#include <string>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <queue>

class Solution {
public:
    std::string alienOrder(std::vector<std::string>& words) {
        std::unordered_map<char, std::unordered_set<char>> adj;
        std::unordered_map<char, int> inDegree;

        for (auto& word : words) {
            for (char ch : word) {
                if (inDegree.find(ch) == inDegree.end()) {
                    inDegree[ch] = 0;
                    adj[ch] = {};
                }
            }
        }

        for (int i = 0; i < static_cast<int>(words.size()) - 1; i++) {
            std::string& w1 = words[i];
            std::string& w2 = words[i + 1];
            int minLen = std::min(w1.size(), w2.size());
            if (w1.size() > w2.size() && w1.substr(0, minLen) == w2.substr(0, minLen)) {
                return "";
            }
            for (int j = 0; j < minLen; j++) {
                if (w1[j] != w2[j]) {
                    if (adj[w1[j]].find(w2[j]) == adj[w1[j]].end()) {
                        adj[w1[j]].insert(w2[j]);
                        inDegree[w2[j]]++;
                    }
                    break;
                }
            }
        }

        std::queue<char> q;
        for (auto& [ch, deg] : inDegree) {
            if (deg == 0) q.push(ch);
        }

        std::string result;
        while (!q.empty()) {
            char ch = q.front();
            q.pop();
            result += ch;
            for (char neighbor : adj[ch]) {
                inDegree[neighbor]--;
                if (inDegree[neighbor] == 0) q.push(neighbor);
            }
        }

        if (result.size() != inDegree.size()) return "";
        return result;
    }
};
