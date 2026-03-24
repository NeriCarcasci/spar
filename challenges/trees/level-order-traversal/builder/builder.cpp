#include <algorithm>
#include <cmath>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <map>
#include <sstream>
#include <stdexcept>
#include <string>
#include <utility>
#include <variant>
#include <vector>

#include "solution.cpp"

struct Value {
    using Array = std::vector<Value>;
    using Object = std::map<std::string, Value>;
    std::variant<std::nullptr_t, bool, double, std::string, Array, Object> data;
    Value() : data(nullptr) {}
    Value(std::nullptr_t) : data(nullptr) {}
    Value(bool v) : data(v) {}
    Value(double v) : data(v) {}
    Value(const std::string& v) : data(v) {}
    Value(const char* v) : data(std::string(v)) {}
    Value(const Array& v) : data(v) {}
    Value(const Object& v) : data(v) {}
    bool isNull() const { return std::holds_alternative<std::nullptr_t>(data); }
    static Value array(std::initializer_list<Value> init) { return Value(Array(init)); }
    static Value object(std::initializer_list<std::pair<std::string, Value>> init) {
        Object obj;
        for (const auto& item : init) {
            obj[item.first] = item.second;
        }
        return Value(obj);
    }
};

double toDoubleValue(const Value& v);
int toInt(const Value& v);
std::string toStringValue(const Value& v);
bool toBoolValue(const Value& v);
std::vector<int> toIntVector(const Value& v);
std::vector<std::vector<int>> toIntMatrix(const Value& v);
std::vector<std::string> toStringVector(const Value& v);
std::vector<std::vector<std::string>> toStringMatrix(const Value& v);
std::vector<char> toCharVector(const Value& v);
std::vector<std::vector<char>> toCharMatrix(const Value& v);

Value toValue(int v) { return Value(static_cast<double>(v)); }
Value toValue(double v) { return Value(v); }
Value toValue(bool v) { return Value(v); }
Value toValue(const std::string& v) { return Value(v); }
Value toValue(const char* v) { return Value(std::string(v)); }
Value toValue(const Value& v) { return v; }

template <typename T>
Value toValue(const std::vector<T>& values) {
    Value::Array out;
    for (const auto& item : values) {
        out.push_back(toValue(item));
    }
    return Value(out);
}

const Value::Array& asArray(const Value& v) {
    static const Value::Array empty;
    if (auto ptr = std::get_if<Value::Array>(&v.data)) {
        return *ptr;
    }
    return empty;
}

Value::Object asObject(const Value& v) {
    if (auto ptr = std::get_if<Value::Object>(&v.data)) {
        return *ptr;
    }
    return Value::Object{};
}

Value getField(const Value::Object& obj, const std::string& key) {
    auto it = obj.find(key);
    if (it == obj.end()) {
        return Value();
    }
    return it->second;
}

bool hasField(const Value::Object& obj, const std::string& key) {
    return obj.find(key) != obj.end();
}

std::string numberToString(double value) {
    if (std::fabs(value - std::round(value)) < 1e-9) {
        long long iv = static_cast<long long>(std::llround(value));
        return std::to_string(iv);
    }
    std::ostringstream oss;
    oss << std::setprecision(15) << value;
    std::string out = oss.str();
    while (!out.empty() && out.back() == '0' && out.find('.') != std::string::npos) {
        out.pop_back();
    }
    if (!out.empty() && out.back() == '.') {
        out.pop_back();
    }
    return out;
}

std::string escapeString(const std::string& input) {
    std::string out;
    out.push_back('"');
    for (char c : input) {
        if (c == '\\' || c == '"') {
            out.push_back('\\');
            out.push_back(c);
        } else if (c == '\n') {
            out += "\\n";
        } else if (c == '\r') {
            out += "\\r";
        } else if (c == '\t') {
            out += "\\t";
        } else {
            out.push_back(c);
        }
    }
    out.push_back('"');
    return out;
}

std::string render(const Value& value);

std::string render(const Value& value) {
    if (std::holds_alternative<std::nullptr_t>(value.data)) {
        return "null";
    }
    if (auto ptr = std::get_if<bool>(&value.data)) {
        return *ptr ? "true" : "false";
    }
    if (auto ptr = std::get_if<double>(&value.data)) {
        return numberToString(*ptr);
    }
    if (auto ptr = std::get_if<std::string>(&value.data)) {
        return escapeString(*ptr);
    }
    if (auto ptr = std::get_if<Value::Array>(&value.data)) {
        std::string out = "[";
        for (size_t i = 0; i < ptr->size(); ++i) {
            if (i > 0) {
                out += ",";
            }
            out += render((*ptr)[i]);
        }
        out += "]";
        return out;
    }
    const auto& obj = std::get<Value::Object>(value.data);
    std::string out = "{";
    bool first = true;
    for (const auto& item : obj) {
        if (!first) {
            out += ",";
        }
        first = false;
        out += escapeString(item.first);
        out += ":";
        out += render(item.second);
    }
    out += "}";
    return out;
}

bool equalsValue(const Value& a, const Value& b) {
    if (a.data.index() != b.data.index()) {
        return false;
    }
    if (std::holds_alternative<std::nullptr_t>(a.data)) {
        return true;
    }
    if (auto pa = std::get_if<bool>(&a.data)) {
        return *pa == std::get<bool>(b.data);
    }
    if (auto pa = std::get_if<double>(&a.data)) {
        return std::fabs(*pa - std::get<double>(b.data)) < 1e-9;
    }
    if (auto pa = std::get_if<std::string>(&a.data)) {
        return *pa == std::get<std::string>(b.data);
    }
    if (auto pa = std::get_if<Value::Array>(&a.data)) {
        const auto& pb = std::get<Value::Array>(b.data);
        if (pa->size() != pb.size()) {
            return false;
        }
        for (size_t i = 0; i < pa->size(); ++i) {
            if (!equalsValue((*pa)[i], pb[i])) {
                return false;
            }
        }
        return true;
    }
    const auto& oa = std::get<Value::Object>(a.data);
    const auto& ob = std::get<Value::Object>(b.data);
    if (oa.size() != ob.size()) {
        return false;
    }
    for (const auto& item : oa) {
        auto it = ob.find(item.first);
        if (it == ob.end()) {
            return false;
        }
        if (!equalsValue(item.second, it->second)) {
            return false;
        }
    }
    return true;
}

Value sortPrimitiveArray(const Value::Array& arr) {
    Value::Array out = arr;
    std::sort(out.begin(), out.end(), [](const Value& left, const Value& right) {
        return render(left) < render(right);
    });
    return Value(out);
}

Value canonical(const Value& v, const std::string& mode) {
    if (mode == "pair_unordered" || mode == "list_unordered" || mode == "strings_unordered") {
        if (!std::holds_alternative<Value::Array>(v.data)) {
            return v;
        }
        return sortPrimitiveArray(std::get<Value::Array>(v.data));
    }
    if (mode == "groups_unordered" || mode == "nested_unordered") {
        if (!std::holds_alternative<Value::Array>(v.data)) {
            return v;
        }
        Value::Array outer;
        for (const Value& item : std::get<Value::Array>(v.data)) {
            if (std::holds_alternative<Value::Array>(item.data)) {
                outer.push_back(sortPrimitiveArray(std::get<Value::Array>(item.data)));
            } else {
                outer.push_back(item);
            }
        }
        std::sort(outer.begin(), outer.end(), [](const Value& left, const Value& right) {
            return render(left) < render(right);
        });
        return Value(outer);
    }
    return v;
}

bool equalFloatSequence(const Value& got, const Value& expected) {
    if (!std::holds_alternative<Value::Array>(got.data) || !std::holds_alternative<Value::Array>(expected.data)) {
        return false;
    }
    const auto& ga = std::get<Value::Array>(got.data);
    const auto& ea = std::get<Value::Array>(expected.data);
    if (ga.size() != ea.size()) {
        return false;
    }
    for (size_t i = 0; i < ga.size(); ++i) {
        if (ga[i].isNull() && ea[i].isNull()) {
            continue;
        }
        if (ga[i].isNull() || ea[i].isNull()) {
            return false;
        }
        if (std::fabs(toDoubleValue(ga[i]) - toDoubleValue(ea[i])) > 1e-9) {
            return false;
        }
    }
    return true;
}

bool equalValues(const Value& got, const Value& expected, const std::string& mode) {
    if (mode == "pair_target_1idx") {
        return false;
    }
    if (mode == "float_sequence") {
        return equalFloatSequence(got, expected);
    }
    Value cg = canonical(got, mode);
    Value ce = canonical(expected, mode);
    if (equalsValue(cg, ce)) {
        return true;
    }
    if (cg.isNull() && std::holds_alternative<Value::Array>(ce.data) && std::get<Value::Array>(ce.data).empty()) {
        return true;
    }
    if (ce.isNull() && std::holds_alternative<Value::Array>(cg.data) && std::get<Value::Array>(cg.data).empty()) {
        return true;
    }
    return false;
}

double toDoubleValue(const Value& v) {
    if (auto ptr = std::get_if<double>(&v.data)) {
        return *ptr;
    }
    if (auto ptr = std::get_if<bool>(&v.data)) {
        return *ptr ? 1.0 : 0.0;
    }
    return 0.0;
}

int toInt(const Value& v) {
    return static_cast<int>(std::llround(toDoubleValue(v)));
}

std::string toStringValue(const Value& v) {
    if (auto ptr = std::get_if<std::string>(&v.data)) {
        return *ptr;
    }
    if (auto ptr = std::get_if<double>(&v.data)) {
        return numberToString(*ptr);
    }
    if (auto ptr = std::get_if<bool>(&v.data)) {
        return *ptr ? "true" : "false";
    }
    return "";
}

bool toBoolValue(const Value& v) {
    if (auto ptr = std::get_if<bool>(&v.data)) {
        return *ptr;
    }
    return false;
}

std::vector<int> toIntVector(const Value& v) {
    std::vector<int> out;
    for (const Value& item : asArray(v)) {
        out.push_back(toInt(item));
    }
    return out;
}

std::vector<std::vector<int>> toIntMatrix(const Value& v) {
    std::vector<std::vector<int>> out;
    for (const Value& row : asArray(v)) {
        out.push_back(toIntVector(row));
    }
    return out;
}

std::vector<std::string> toStringVector(const Value& v) {
    std::vector<std::string> out;
    for (const Value& item : asArray(v)) {
        out.push_back(toStringValue(item));
    }
    return out;
}

std::vector<std::vector<std::string>> toStringMatrix(const Value& v) {
    std::vector<std::vector<std::string>> out;
    for (const Value& row : asArray(v)) {
        out.push_back(toStringVector(row));
    }
    return out;
}

std::vector<char> toCharVector(const Value& v) {
    std::vector<char> out;
    for (const Value& item : asArray(v)) {
        std::string s = toStringValue(item);
        out.push_back(s.empty() ? '\0' : s[0]);
    }
    return out;
}

std::vector<std::vector<char>> toCharMatrix(const Value& v) {
    std::vector<std::vector<char>> out;
    for (const Value& row : asArray(v)) {
        out.push_back(toCharVector(row));
    }
    return out;
}

Value getArg(const Value::Object& testObj, const std::vector<std::string>& argKeys, size_t idx) {
    auto itInput = testObj.find("input");
    if (itInput != testObj.end()) {
        const Value& input = itInput->second;
        if (std::holds_alternative<Value::Object>(input.data)) {
            Value::Object inputObj = std::get<Value::Object>(input.data);
            if (idx < argKeys.size()) {
                return getField(inputObj, argKeys[idx]);
            }
            return Value();
        }
        if (std::holds_alternative<Value::Array>(input.data)) {
            const auto& arr = std::get<Value::Array>(input.data);
            if (argKeys.size() > 1) {
                if (idx == 0) {
                    return input;
                }
                if (idx < argKeys.size()) {
                    auto it = testObj.find(argKeys[idx]);
                    if (it != testObj.end()) {
                        return it->second;
                    }
                }
                if (idx < arr.size()) {
                    return arr[idx];
                }
                return Value();
            }
            return input;
        }
        return input;
    }
    if (idx < argKeys.size()) {
        auto it = testObj.find(argKeys[idx]);
        if (it != testObj.end()) {
            return it->second;
        }
    }
    return Value();
}



TreeNode* buildTree(const Value& v) {
    const Value::Array& raw = asArray(v);
    if (raw.empty() || raw[0].isNull()) {
        return nullptr;
    }
    std::vector<TreeNode*> nodes(raw.size(), nullptr);
    for (size_t i = 0; i < raw.size(); ++i) {
        if (raw[i].isNull()) {
            continue;
        }
        nodes[i] = new TreeNode(toInt(raw[i]));
    }
    size_t idx = 1;
    for (TreeNode* node : nodes) {
        if (node == nullptr) {
            continue;
        }
        if (idx < nodes.size()) {
            node->left = nodes[idx++];
        }
        if (idx < nodes.size()) {
            node->right = nodes[idx++];
        }
    }
    return nodes[0];
}

Value treeToValue(TreeNode* root) {
    if (root == nullptr) {
        return Value::array({});
    }
    std::vector<TreeNode*> queue = {root};
    Value::Array out;
    for (size_t i = 0; i < queue.size(); ++i) {
        TreeNode* node = queue[i];
        if (node == nullptr) {
            out.push_back(Value());
            continue;
        }
        out.push_back(toValue(node->val));
        queue.push_back(node->left);
        queue.push_back(node->right);
    }
    while (!out.empty() && out.back().isNull()) {
        out.pop_back();
    }
    return Value(out);
}


Node* buildGraph(const Value& v) {
    std::vector<std::vector<int>> adj = toIntMatrix(v);
    if (adj.empty()) {
        return nullptr;
    }
    std::vector<Node*> nodes;
    nodes.reserve(adj.size());
    for (size_t i = 0; i < adj.size(); ++i) {
        nodes.push_back(new Node(static_cast<int>(i + 1)));
    }
    for (size_t i = 0; i < adj.size(); ++i) {
        std::vector<Node*> neighbors;
        for (int n : adj[i]) {
            if (n >= 1 && static_cast<size_t>(n) <= nodes.size()) {
                neighbors.push_back(nodes[n - 1]);
            }
        }
        nodes[i]->neighbors = neighbors;
    }
    return nodes[0];
}

Value graphToValue(Node* root) {
    if (root == nullptr) {
        return Value::array({});
    }
    std::map<int, Node*> seen;
    std::vector<Node*> queue = {root};
    for (size_t i = 0; i < queue.size(); ++i) {
        Node* node = queue[i];
        if (node == nullptr) {
            continue;
        }
        if (seen.find(node->val) != seen.end()) {
            continue;
        }
        seen[node->val] = node;
        for (Node* nei : node->neighbors) {
            queue.push_back(nei);
        }
    }
    int maxId = 0;
    for (const auto& item : seen) {
        if (item.first > maxId) {
            maxId = item.first;
        }
    }
    Value::Array out(maxId, Value::array({}));
    for (const auto& item : seen) {
        std::vector<int> row;
        for (Node* nei : item.second->neighbors) {
            row.push_back(nei->val);
        }
        std::sort(row.begin(), row.end());
        out[item.first - 1] = toValue(row);
    }
    return Value(out);
}


int main(int argc, char** argv) {
    std::string testsPath = "../tests.yaml";
    std::string challengePath = "../challenge.yaml";
    if (argc > 2) {
        testsPath = argv[2];
    }
    if (argc > 3) {
        challengePath = argv[3];
    }
    std::ifstream testsFile(testsPath);
    std::ifstream challengeFile(challengePath);
    (void)testsFile;
    (void)challengeFile;

    std::vector<Value> tests = {
        Value::object({{"input", Value::object({{"root", Value::array({Value(3.0), Value(9.0), Value(20.0), Value(), Value(), Value(15.0), Value(7.0)})}})}, {"expected", Value::array({Value::array({Value(3.0)}), Value::array({Value(9.0), Value(20.0)}), Value::array({Value(15.0), Value(7.0)})})}, {"visible", Value(true)}}),
        Value::object({{"input", Value::object({{"root", Value::array({Value(1.0)})}})}, {"expected", Value::array({Value::array({Value(1.0)})})}, {"visible", Value(true)}}),
        Value::object({{"input", Value::object({{"root", Value::array({})}})}, {"expected", Value::array({})}, {"visible", Value(false)}}),
        Value::object({{"input", Value::object({{"root", Value::array({Value(1.0), Value(2.0), Value(3.0), Value(4.0), Value(5.0)})}})}, {"expected", Value::array({Value::array({Value(1.0)}), Value::array({Value(2.0), Value(3.0)}), Value::array({Value(4.0), Value(5.0)})})}, {"visible", Value(false)}}),
        Value::object({{"input", Value::object({{"root", Value::array({Value(1.0), Value(), Value(2.0), Value(), Value(3.0)})}})}, {"expected", Value::array({Value::array({Value(1.0)}), Value::array({Value(2.0)}), Value::array({Value(3.0)})})}, {"visible", Value(false)}})
    };
    std::vector<std::string> argKeys = {"root"};
    bool failed = false;

    for (size_t i = 0; i < tests.size(); ++i) {
        Value::Object testObj = asObject(tests[i]);
        Value gotValue;
        Solution sol;
        Value raw0 = getArg(testObj, argKeys, 0);
        auto arg0 = buildTree(raw0);
        auto result = sol.levelOrder(arg0);
        gotValue = toValue(result);

        Value expected = getField(testObj, "expected");
        bool hasExpectedLength = hasField(testObj, "expected_length");
        bool passed = false;
        if (hasExpectedLength) {
            Value expectedLength = getField(testObj, "expected_length");
            if (std::holds_alternative<Value::Array>(gotValue.data)) {
                passed = static_cast<int>(std::get<Value::Array>(gotValue.data).size()) == toInt(expectedLength);
            } else {
                passed = false;
            }
        } else if ("exact" == std::string("pair_target_1idx")) {
            Value::Object input = asObject(getField(testObj, "input"));
            Value::Array numbers = asArray(getField(input, "numbers"));
            int target = toInt(getField(input, "target"));
            if (std::holds_alternative<Value::Array>(gotValue.data)) {
                const Value::Array& pair = std::get<Value::Array>(gotValue.data);
                if (pair.size() == 2) {
                    int i1 = toInt(pair[0]) - 1;
                    int i2 = toInt(pair[1]) - 1;
                    if (i1 >= 0 && i2 >= 0 && i1 < static_cast<int>(numbers.size()) && i2 < static_cast<int>(numbers.size()) && i1 != i2) {
                        passed = toInt(numbers[i1]) + toInt(numbers[i2]) == target;
                    } else {
                        passed = false;
                    }
                } else {
                    passed = false;
                }
            } else {
                passed = false;
            }
        } else {
            passed = equalValues(gotValue, expected, "exact");
        }

        int idx = static_cast<int>(i + 1);
        if (passed) {
            std::cout << "PASS " << idx << "\n";
        } else {
            failed = true;
            if (hasExpectedLength) {
                std::cout << "FAIL " << idx << " got=" << render(gotValue) << " expected=" << render(getField(testObj, "expected_length")) << "\n";
            } else {
                std::cout << "FAIL " << idx << " got=" << render(gotValue) << " expected=" << render(expected) << "\n";
            }
        }
    }

    if (failed) {
        return 1;
    }
    return 0;
}
