#include <vector>
#include <queue>

struct ListNode {
    int val;
    ListNode* next;
    ListNode(int x) : val(x), next(nullptr) {}
};

class Solution {
public:
    ListNode* mergeKLists(std::vector<ListNode*>& lists) {
        auto cmp = [](ListNode* a, ListNode* b) { return a->val > b->val; };
        std::priority_queue<ListNode*, std::vector<ListNode*>, decltype(cmp)> pq(cmp);
        for (auto* node : lists) {
            if (node) pq.push(node);
        }
        ListNode dummy(0);
        ListNode* curr = &dummy;
        while (!pq.empty()) {
            ListNode* node = pq.top(); pq.pop();
            curr->next = node;
            curr = curr->next;
            if (node->next) pq.push(node->next);
        }
        return dummy.next;
    }
};
