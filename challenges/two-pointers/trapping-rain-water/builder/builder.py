import importlib.util
import json
import pathlib
import sys

TESTS = json.loads('[{"input":{"height":[0,1,0,2,1,0,1,3,2,1,2,1]},"expected":6,"visible":true},{"input":{"height":[4,2,0,3,2,5]},"expected":9,"visible":true},{"input":{"height":[1]},"expected":0,"visible":false},{"input":{"height":[1,2,3,4,5]},"expected":0,"visible":false},{"input":{"height":[5,4,3,2,1]},"expected":0,"visible":false},{"input":{"height":[3,0,0,0,3]},"expected":9,"visible":false},{"input":{"height":[0,0,0]},"expected":0,"visible":false}]')
COMPARE_MODE = 'exact'


def resolve_callable(module, names):
    for name in names:
        fn = getattr(module, name, None)
        if callable(fn):
            return fn
    raise RuntimeError("missing function")


def resolve_class(module, name, fallback):
    cls = getattr(module, name, None)
    if cls is not None:
        return cls
    return fallback


def list_node_fallback(val=0, next=None):
    class ListNode:
        def __init__(self, val=0, next=None):
            self.val = val
            self.next = next
    return ListNode(val, next)


def tree_node_fallback(val=0, left=None, right=None):
    class TreeNode:
        def __init__(self, val=0, left=None, right=None):
            self.val = val
            self.left = left
            self.right = right
    return TreeNode(val, left, right)


def graph_node_fallback(val=0, neighbors=None):
    class Node:
        def __init__(self, val=0, neighbors=None):
            self.val = val
            self.neighbors = neighbors if neighbors is not None else []
    return Node(val, neighbors)


def array_to_linked_list(values, node_cls):
    head = None
    tail = None
    for value in values:
        node = node_cls(value)
        if head is None:
            head = node
            tail = node
        else:
            tail.next = node
            tail = node
    return head


def linked_list_to_array(head):
    out = []
    node = head
    seen = 0
    while node is not None and seen < 100000:
        out.append(node.val)
        node = node.next
        seen += 1
    return out


def array_to_tree(values, node_cls):
    if not values:
        return None
    if values[0] is None:
        return None
    nodes = [None if v is None else node_cls(v) for v in values]
    idx = 1
    for node in nodes:
        if node is None:
            continue
        if idx < len(nodes):
            node.left = nodes[idx]
            idx += 1
        if idx < len(nodes):
            node.right = nodes[idx]
            idx += 1
    return nodes[0]


def tree_to_array(root):
    if root is None:
        return []
    out = []
    queue = [root]
    i = 0
    while i < len(queue):
        node = queue[i]
        i += 1
        if node is None:
            out.append(None)
            continue
        out.append(node.val)
        queue.append(node.left)
        queue.append(node.right)
    while out and out[-1] is None:
        out.pop()
    return out


def find_tree_node(root, value):
    if root is None:
        return None
    queue = [root]
    i = 0
    while i < len(queue):
        node = queue[i]
        i += 1
        if node is None:
            continue
        if node.val == value:
            return node
        queue.append(node.left)
        queue.append(node.right)
    return None


def adjlist_to_graph(adj_list, node_cls):
    if not adj_list:
        return None
    nodes = [node_cls(i + 1) for i in range(len(adj_list))]
    for idx, neighbors in enumerate(adj_list):
        nodes[idx].neighbors = [nodes[n - 1] for n in neighbors]
    return nodes[0]


def graph_to_adjlist(root):
    if root is None:
        return []
    seen = {}
    queue = [root]
    order = []
    while queue:
        node = queue.pop(0)
        if node.val in seen:
            continue
        seen[node.val] = node
        order.append(node.val)
        for nei in node.neighbors:
            queue.append(nei)
    max_id = max(order) if order else 0
    out = [[] for _ in range(max_id)]
    for value, node in seen.items():
        out[value - 1] = sorted([n.val for n in node.neighbors])
    return out


def render(value):
    return json.dumps(value, separators=(",", ":"))


def canonical(value, mode):
    if mode == "pair_unordered":
        if isinstance(value, list):
            return sorted(value)
        return value
    if mode == "list_unordered":
        if isinstance(value, list):
            return sorted(value)
        return value
    if mode == "strings_unordered":
        if isinstance(value, list):
            return sorted(value)
        return value
    if mode == "groups_unordered":
        if isinstance(value, list):
            groups = []
            for group in value:
                if isinstance(group, list):
                    groups.append(sorted(group))
                else:
                    groups.append(group)
            return sorted(groups)
        return value
    if mode == "nested_unordered":
        if isinstance(value, list):
            normalized = []
            for item in value:
                if isinstance(item, list):
                    normalized.append(sorted(item))
                else:
                    normalized.append(item)
            return sorted(normalized)
        return value
    return value


def equal_values(got, expected, mode):
    if mode == "pair_target_1idx":
        return False
    if mode == "float_sequence":
        if not isinstance(got, list) or not isinstance(expected, list) or len(got) != len(expected):
            return False
        for a, b in zip(got, expected):
            if a is None and b is None:
                continue
            if a is None or b is None:
                return False
            if abs(float(a) - float(b)) > 1e-9:
                return False
        return True
    return canonical(got, mode) == canonical(expected, mode)


def load_module(solution_path):
    spec = importlib.util.spec_from_file_location("user_solution", solution_path)
    if spec is None or spec.loader is None:
        raise RuntimeError("unable to load solution")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def run(solution_path, tests_path, challenge_path):
    _ = pathlib.Path(tests_path).read_text(encoding="utf-8")
    _ = pathlib.Path(challenge_path).read_text(encoding="utf-8")
    module = load_module(solution_path)
    failures = 0
    for idx, case in enumerate(TESTS, start=1):

        fn = resolve_callable(module, ['trap'])
        case_input = case.get("input")
        args = []
        if isinstance(case_input, dict):
            for key in ['height']:
                args.append(case_input.get(key))
        elif isinstance(case_input, list) and len(['height']) > 1:
            args.append(case_input)
            for idx in range(1, len(['height'])):
                key = ['height'][idx]
                if key in case:
                    args.append(case.get(key))
                else:
                    args.append(case_input[idx] if idx < len(case_input) else None)
        elif case_input is not None:
            args = [case_input]
        else:
            for key in ['height']:
                args.append(case.get(key))
        got_value = fn(*args)

        expected = case.get("expected")
        expected_length = case.get("expected_length")
        passed = False
        if COMPARE_MODE == "pair_target_1idx":
            inp = case.get("input", {})
            numbers = inp.get("numbers", []) if isinstance(inp, dict) else []
            target = inp.get("target") if isinstance(inp, dict) else None
            if isinstance(got_value, list) and len(got_value) == 2 and target is not None:
                i = got_value[0] - 1
                j = got_value[1] - 1
                passed = 0 <= i < len(numbers) and 0 <= j < len(numbers) and i != j and (numbers[i] + numbers[j] == target)
            else:
                passed = False
        elif expected_length is not None:
            try:
                passed = len(got_value) == expected_length
            except Exception:
                passed = False
        else:
            passed = equal_values(got_value, expected, COMPARE_MODE)
        if passed:
            print(f"PASS {idx}")
        else:
            failures += 1
            if expected_length is not None:
                print(f"FAIL {idx} got={render(got_value)} expected={render(expected_length)}")
            else:
                print(f"FAIL {idx} got={render(got_value)} expected={render(expected)}")
    return failures == 0


def main():
    builder_dir = pathlib.Path(__file__).resolve().parent
    solution_path = pathlib.Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else (builder_dir.parent / "setup" / "python.py")
    tests_path = pathlib.Path(sys.argv[2]).resolve() if len(sys.argv) > 2 else (builder_dir.parent / "tests.yaml")
    challenge_path = pathlib.Path(sys.argv[3]).resolve() if len(sys.argv) > 3 else (builder_dir.parent / "challenge.yaml")
    ok = run(solution_path, tests_path, challenge_path)
    raise SystemExit(0 if ok else 1)


if __name__ == "__main__":
    main()
