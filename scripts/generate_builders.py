import json
import os
import re
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
CHALLENGES_DIR = ROOT / "challenges"

ORDER_INSENSITIVE_NESTED = {
    "two-pointers/three-sum",
    "backtracking/combination-sum",
    "backtracking/permutations",
    "backtracking/subsets",
    "backtracking/n-queens",
    "graphs/pacific-atlantic",
    "heap/merge-k-sorted-alt",
}
ORDER_INSENSITIVE_STRINGS = {
    "tries/word-search-ii",
}
ORDER_INSENSITIVE_GROUPS = {
    "arrays-and-hashing/group-anagrams",
}
ORDER_INSENSITIVE_SCALAR_LIST = {
    "arrays-and-hashing/top-k-frequent",
}
PAIR_ORDER_INSENSITIVE = {
    "arrays-and-hashing/two-sum",
}
FLOAT_SEQUENCE = {
    "heap/find-median-stream",
}
PAIR_TARGET_1INDEXED = {
    "two-pointers/two-sum-sorted",
}

LINKED_LIST_CHALLENGES = {
    "linked-lists/reverse-linked-list",
    "linked-lists/merge-two-sorted",
    "linked-lists/merge-k-sorted",
    "linked-lists/remove-nth-from-end",
    "linked-lists/reorder-list",
    "linked-lists/linked-list-cycle",
}
TREE_CHALLENGES = {
    "trees/binary-tree-max-path",
    "trees/build-tree-preorder",
    "trees/invert-binary-tree",
    "trees/kth-smallest-bst",
    "trees/level-order-traversal",
    "trees/lowest-common-ancestor",
    "trees/max-depth-tree",
    "trees/same-tree",
    "trees/serialize-deserialize",
    "trees/subtree-of-another",
    "trees/validate-bst",
}

SPECIAL_KIND = {
    "arrays-and-hashing/encode-decode-strings": "encode_decode",
    "trees/serialize-deserialize": "serialize_deserialize",
    "heap/find-median-stream": "class_median_finder",
    "heap/kth-largest-stream": "class_kth_largest",
    "stack/min-stack": "class_min_stack",
    "tries/implement-trie": "class_trie",
    "linked-lists/linked-list-cycle": "linked_list_cycle",
    "linked-lists/merge-k-sorted": "linked_list_merge_k",
    "linked-lists/merge-two-sorted": "linked_list_merge_two",
    "linked-lists/remove-nth-from-end": "linked_list_remove_nth",
    "linked-lists/reorder-list": "linked_list_reorder",
    "linked-lists/reverse-linked-list": "linked_list_reverse",
    "trees/lowest-common-ancestor": "tree_lca",
    "graphs/clone-graph": "clone_graph",
}


def discover_challenges():
    result = []
    for tests_path in CHALLENGES_DIR.rglob("tests.yaml"):
        rel = tests_path.parent.relative_to(CHALLENGES_DIR).as_posix()
        result.append(rel)
    return sorted(result)


def load_yaml(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def normalize_tests(data):
    if isinstance(data, dict) and "tests" in data:
        visible = []
        hidden = []
        for case in data.get("tests", []):
            if case.get("visible"):
                visible.append(case)
            else:
                hidden.append(case)
        return visible + hidden
    visible = data.get("visible", []) if isinstance(data, dict) else []
    hidden = data.get("hidden", []) if isinstance(data, dict) else []
    return list(visible) + list(hidden)


def parse_py_signature(path: Path):
    text = path.read_text(encoding="utf-8")
    funcs = re.findall(r"^def\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)", text, re.M)
    classes = re.findall(r"^class\s+([A-Za-z_][A-Za-z0-9_]*)\b", text, re.M)
    out = {}
    for name, params in funcs:
        parts = []
        for p in params.split(","):
            p = p.strip()
            if not p:
                continue
            p = p.split(":", 1)[0].strip()
            p = p.split("=", 1)[0].strip()
            parts.append(p)
        out[name] = parts
    return out, classes


def arg_keys_from_tests(tests):
    for case in tests:
        if not isinstance(case, dict):
            continue
        if "input" in case and isinstance(case["input"], dict):
            return list(case["input"].keys())
    return []


def compare_mode(challenge):
    if challenge in PAIR_TARGET_1INDEXED:
        return "pair_target_1idx"
    if challenge in PAIR_ORDER_INSENSITIVE:
        return "pair_unordered"
    if challenge in ORDER_INSENSITIVE_SCALAR_LIST:
        return "list_unordered"
    if challenge in ORDER_INSENSITIVE_GROUPS:
        return "groups_unordered"
    if challenge in ORDER_INSENSITIVE_STRINGS:
        return "strings_unordered"
    if challenge in ORDER_INSENSITIVE_NESTED:
        return "nested_unordered"
    if challenge in FLOAT_SEQUENCE:
        return "float_sequence"
    return "exact"


def python_kind(challenge):
    if challenge in SPECIAL_KIND:
        return SPECIAL_KIND[challenge]
    if challenge in TREE_CHALLENGES:
        return "tree_plain"
    return "plain"


def py_functions_for_challenge(funcs):
    names = list(funcs.keys())
    return names


def parse_js_signature(path: Path):
    text = path.read_text(encoding="utf-8")
    funcs = re.findall(r"function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(", text)
    classes = re.findall(r"class\s+([A-Za-z_][A-Za-z0-9_]*)\b", text)
    return funcs, classes


def parse_go_signature(path: Path):
    text = path.read_text(encoding="utf-8")
    funcs = {}
    methods = {}
    for m in re.finditer(r"^func\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*([^{\n]*)\{", text, re.M):
        name = m.group(1)
        params = parse_go_params(m.group(2))
        ret = m.group(3).strip()
        funcs[name] = {"params": params, "ret": ret}
    for m in re.finditer(r"^func\s*\([^)]*\)\s*([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*([^{\n]*)\{", text, re.M):
        name = m.group(1)
        params = parse_go_params(m.group(2))
        ret = m.group(3).strip()
        methods[name] = {"params": params, "ret": ret}
    return funcs, methods


def parse_go_params(raw: str):
    raw = raw.strip()
    if not raw:
        return []
    parts = []
    segments = [p.strip() for p in raw.split(",")]
    for seg in segments:
        if not seg:
            continue
        tokens = seg.split()
        if len(tokens) < 2:
            continue
        ptype = tokens[-1].strip()
        names = " ".join(tokens[:-1]).split(",")
        for name in names:
            name = name.strip()
            if name:
                parts.append((name, ptype))
    return parts


def parse_cpp_signature(path: Path):
    text = path.read_text(encoding="utf-8")
    solution_methods = {}
    class_methods = {}
    classes = {}

    for m in re.finditer(r"class\s+([A-Za-z_][A-Za-z0-9_]*)\s*\{", text):
        cname = m.group(1)
        i = m.end()
        depth = 1
        while i < len(text) and depth > 0:
            ch = text[i]
            if ch == "{":
                depth += 1
            elif ch == "}":
                depth -= 1
            i += 1
        block = text[m.end() : i - 1]
        methods = []
        for mm in re.finditer(r"([A-Za-z_][A-Za-z0-9_:<>\s*&]+?)\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*\{", block):
            ret = mm.group(1).replace("public:", "").replace("private:", "").replace("protected:", "").strip()
            name = mm.group(2).strip()
            params = parse_cpp_params(mm.group(3))
            methods.append({"name": name, "ret": ret, "params": params})
        classes[cname] = methods
    if "Solution" in classes:
        for method in classes["Solution"]:
            solution_methods[method["name"]] = method
    for cname, methods_list in classes.items():
        if cname == "Solution":
            continue
        class_methods[cname] = {m["name"]: m for m in methods_list}
    return solution_methods, class_methods


def parse_cpp_params(raw: str):
    raw = raw.strip()
    if not raw:
        return []
    parts = []
    segments = [seg.strip() for seg in raw.split(",")]
    for seg in segments:
        if not seg:
            continue
        tokens = seg.split()
        if len(tokens) < 2:
            continue
        name = tokens[-1].strip()
        ptype = " ".join(tokens[:-1]).strip()
        parts.append((name, ptype))
    return parts


def split_params_with_depth(raw: str):
    parts = []
    current = []
    depth = 0
    for ch in raw:
        if ch == "," and depth == 0:
            part = "".join(current).strip()
            if part:
                parts.append(part)
            current = []
            continue
        if ch in "<([{" :
            depth += 1
        elif ch in ">)]}":
            depth = max(0, depth - 1)
        current.append(ch)
    part = "".join(current).strip()
    if part:
        parts.append(part)
    return parts


def parse_rust_signature(path: Path):
    text = path.read_text(encoding="utf-8")
    funcs = {}
    impl_methods = {}
    structs = []

    for m in re.finditer(r"pub\s+struct\s+([A-Za-z_][A-Za-z0-9_]*)", text):
        structs.append(m.group(1))

    for m in re.finditer(r"^pub\s+fn\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*(?:->\s*([^{]+))?\{", text, re.M):
        name = m.group(1)
        params = parse_rust_params(m.group(2))
        ret = (m.group(3) or "").strip()
        funcs[name] = {"params": params, "ret": ret}

    for m in re.finditer(r"impl\s+([A-Za-z_][A-Za-z0-9_]*)\s*\{", text):
        sname = m.group(1)
        i = m.end()
        depth = 1
        while i < len(text) and depth > 0:
            ch = text[i]
            if ch == "{":
                depth += 1
            elif ch == "}":
                depth -= 1
            i += 1
        block = text[m.end() : i - 1]
        methods = {}
        for mm in re.finditer(r"pub\s+fn\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(([^)]*)\)\s*(?:->\s*([^{]+))?\{", block):
            mname = mm.group(1)
            params = parse_rust_params(mm.group(2))
            ret = (mm.group(3) or "").strip()
            methods[mname] = {"params": params, "ret": ret}
        impl_methods[sname] = methods

    return funcs, impl_methods, structs


def parse_rust_params(raw: str):
    raw = raw.strip()
    if not raw:
        return []
    parts = []
    for seg in split_params_with_depth(raw):
        if seg in ("&self", "&mut self", "self"):
            parts.append(("self", seg))
            continue
        if ":" not in seg:
            continue
        name, ptype = seg.split(":", 1)
        parts.append((name.strip(), ptype.strip()))
    return parts


def py_builder(challenge, tests, funcs, classes, arg_keys):
    kind = python_kind(challenge)
    mode = compare_mode(challenge)
    function_names = py_functions_for_challenge(funcs)
    tests_json = json.dumps(tests, separators=(",", ":"), ensure_ascii=False)

    call_block = ""
    if kind == "encode_decode":
        call_block = """
        encode = resolve_callable(module, [\"encode\"])
        decode = resolve_callable(module, [\"decode\"])
        inp = case.get(\"input\", {})
        raw = inp.get(\"strs\", []) if isinstance(inp, dict) else inp
        got_value = decode(encode(raw))
"""
    elif kind == "serialize_deserialize":
        call_block = """
        serialize = resolve_callable(module, [\"serialize\"])
        deserialize = resolve_callable(module, [\"deserialize\"])
        node_cls = resolve_class(module, \"TreeNode\", tree_node_fallback)
        inp = case.get(\"input\", {})
        raw = inp.get(\"root\", []) if isinstance(inp, dict) else inp
        root = array_to_tree(raw, node_cls)
        got_value = tree_to_array(deserialize(serialize(root)))
"""
    elif kind == "class_median_finder":
        call_block = """
        cls = resolve_class(module, \"MedianFinder\", None)
        if cls is None:
            raise RuntimeError(\"missing MedianFinder class\")
        inp = case.get(\"input\", {})
        operations = inp.get(\"operations\", [])
        values = inp.get(\"values\", [])
        obj = cls()
        out = []
        for op, val in zip(operations, values):
            if op == \"add_num\":
                obj.add_num(val)
                out.append(None)
            elif op == \"find_median\":
                out.append(obj.find_median())
            else:
                raise RuntimeError(f\"unknown op {op}\")
        got_value = out
"""
    elif kind == "class_kth_largest":
        call_block = """
        cls = resolve_class(module, \"KthLargest\", None)
        if cls is None:
            raise RuntimeError(\"missing KthLargest class\")
        inp = case.get(\"input\", {})
        obj = cls(inp.get(\"k\"), inp.get(\"nums\", []))
        operations = inp.get(\"operations\", [])
        values = inp.get(\"values\", [])
        out = []
        for op, val in zip(operations, values):
            if op != \"add\":
                raise RuntimeError(f\"unknown op {op}\")
            out.append(obj.add(val))
        got_value = out
"""
    elif kind == "class_min_stack":
        call_block = """
        cls = resolve_class(module, \"MinStack\", None)
        if cls is None:
            raise RuntimeError(\"missing MinStack class\")
        inp = case.get(\"input\", {})
        operations = inp.get(\"operations\", [])
        values = inp.get(\"values\", [])
        obj = cls()
        out = []
        for op, val in zip(operations, values):
            if op == \"push\":
                obj.push(val)
                out.append(None)
            elif op == \"pop\":
                obj.pop()
                out.append(None)
            elif op == \"top\":
                out.append(obj.top())
            elif op == \"get_min\":
                out.append(obj.get_min())
            else:
                raise RuntimeError(f\"unknown op {op}\")
        got_value = out
"""
    elif kind == "class_trie":
        call_block = """
        cls = resolve_class(module, \"Trie\", None)
        if cls is None:
            raise RuntimeError(\"missing Trie class\")
        inp = case.get(\"input\", {})
        operations = inp.get(\"operations\", [])
        values = inp.get(\"values\", [])
        obj = cls()
        out = []
        for op, val in zip(operations, values):
            if op == \"insert\":
                obj.insert(val)
                out.append(None)
            elif op == \"search\":
                out.append(obj.search(val))
            elif op == \"starts_with\":
                out.append(obj.starts_with(val))
            else:
                raise RuntimeError(f\"unknown op {op}\")
        got_value = out
"""
    elif kind == "linked_list_cycle":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        inp = case.get(\"input\", {})
        values = inp.get(\"list\", [])
        pos = inp.get(\"pos\", -1)
        head = array_to_linked_list(values, node_cls)
        if pos is not None and pos >= 0:
            tail = head
            target = None
            idx = 0
            while tail is not None and tail.next is not None:
                if idx == pos:
                    target = tail
                tail = tail.next
                idx += 1
            if tail is not None:
                if idx == pos:
                    target = tail
                if target is not None:
                    tail.next = target
        got_value = fn(head)
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "linked_list_merge_k":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        raw = case.get(\"input\", [])
        lists = [array_to_linked_list(arr, node_cls) for arr in raw]
        got_value = linked_list_to_array(fn(lists))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "linked_list_merge_two":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        raw = case.get(\"input\", [])
        left = array_to_linked_list(raw[0] if len(raw) > 0 else [], node_cls)
        right = array_to_linked_list(raw[1] if len(raw) > 1 else [], node_cls)
        got_value = linked_list_to_array(fn(left, right))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "linked_list_remove_nth":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        inp = case.get(\"input\", {})
        head = array_to_linked_list(inp.get(\"list\", []), node_cls)
        got_value = linked_list_to_array(fn(head, inp.get(\"n\")))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "linked_list_reorder":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        head = array_to_linked_list(case.get(\"input\", []), node_cls)
        result = fn(head)
        got_value = linked_list_to_array(head if result is None else result)
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "linked_list_reverse":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"ListNode\", list_node_fallback)
        head = array_to_linked_list(case.get(\"input\", []), node_cls)
        got_value = linked_list_to_array(fn(head))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "tree_lca":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"TreeNode\", tree_node_fallback)
        inp = case.get(\"input\", {})
        root = array_to_tree(inp.get(\"root\", []), node_cls)
        got_value = fn(root, inp.get(\"p\"), inp.get(\"q\"))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "clone_graph":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"Node\", graph_node_fallback)
        inp = case.get(\"input\", {})
        root = adjlist_to_graph(inp.get(\"adjList\", []), node_cls)
        got_value = graph_to_adjlist(fn(root))
""".replace("__FUNCTION_NAMES__", repr(function_names))
    elif kind == "tree_plain":
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        node_cls = resolve_class(module, \"TreeNode\", tree_node_fallback)
        inp = case.get(\"input\", {})
        args = []
        for key in __ARG_KEYS__:
            value = inp.get(key)
            if key in (\"root\", \"subRoot\", \"p\", \"q\"):
                value = array_to_tree(value, node_cls)
            args.append(value)
        out = fn(*args)
        if __CHALLENGE__ in (\"trees/build-tree-preorder\", \"trees/invert-binary-tree\"):
            got_value = tree_to_array(out)
        else:
            got_value = out
""".replace("__FUNCTION_NAMES__", repr(function_names)).replace("__ARG_KEYS__", repr(arg_keys)).replace("__CHALLENGE__", repr(challenge))
    else:
        call_block = """
        fn = resolve_callable(module, __FUNCTION_NAMES__)
        case_input = case.get(\"input\")
        args = []
        if isinstance(case_input, dict):
            for key in __ARG_KEYS__:
                args.append(case_input.get(key))
        elif isinstance(case_input, list) and len(__ARG_KEYS__) > 1:
            args.append(case_input)
            for idx in range(1, len(__ARG_KEYS__)):
                key = __ARG_KEYS__[idx]
                if key in case:
                    args.append(case.get(key))
                else:
                    args.append(case_input[idx] if idx < len(case_input) else None)
        elif case_input is not None:
            args = [case_input]
        else:
            for key in __ARG_KEYS__:
                args.append(case.get(key))
        got_value = fn(*args)
""".replace("__FUNCTION_NAMES__", repr(function_names)).replace("__ARG_KEYS__", repr(arg_keys))

    content = f'''import importlib.util
import json
import pathlib
import sys

TESTS = json.loads({tests_json!r})
COMPARE_MODE = {mode!r}


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
    seen = {{}}
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
{call_block}
        expected = case.get("expected")
        expected_length = case.get("expected_length")
        passed = False
        if COMPARE_MODE == "pair_target_1idx":
            inp = case.get("input", {{}})
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
            print(f"PASS {{idx}}")
        else:
            failures += 1
            if expected_length is not None:
                print(f"FAIL {{idx}} got={{render(got_value)}} expected={{render(expected_length)}}")
            else:
                print(f"FAIL {{idx}} got={{render(got_value)}} expected={{render(expected)}}")
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
'''
    return content


def js_builder(challenge, tests, funcs, classes, arg_keys):
    kind = python_kind(challenge)
    mode = compare_mode(challenge)
    tests_json = json.dumps(tests, separators=(",", ":"), ensure_ascii=False)
    function_names = funcs or []

    if kind == "encode_decode":
        call_block = """
    const encode = resolveFunction(solution, ["encode"]);
    const decode = resolveFunction(solution, ["decode"]);
    const inp = test.input || {};
    const raw = isObject(inp) ? inp.strs || [] : inp;
    gotValue = decode(encode(raw));
"""
    elif kind == "serialize_deserialize":
        call_block = """
    const serialize = resolveFunction(solution, ["serialize"]);
    const deserialize = resolveFunction(solution, ["deserialize"]);
    const TreeNodeClass = resolveClass(solution, "TreeNode", TreeNodeFallback);
    const inp = test.input || {};
    const raw = isObject(inp) ? inp.root || [] : inp;
    const root = arrayToTree(raw, TreeNodeClass);
    gotValue = treeToArray(deserialize(serialize(root)));
"""
    elif kind == "class_median_finder":
        call_block = """
    const MedianFinderClass = resolveClass(solution, "MedianFinder", null);
    if (!MedianFinderClass) {
        throw new Error("missing MedianFinder class");
    }
    const inp = test.input || {};
    const operations = inp.operations || [];
    const values = inp.values || [];
    const obj = new MedianFinderClass();
    const out = [];
    for (let i = 0; i < operations.length; i += 1) {
        const op = operations[i];
        const val = values[i];
        if (op === "add_num") {
            obj.addNum(val);
            out.push(null);
        } else if (op === "find_median") {
            out.push(obj.findMedian());
        } else {
            throw new Error(`unknown op ${op}`);
        }
    }
    gotValue = out;
"""
    elif kind == "class_kth_largest":
        call_block = """
    const KthLargestClass = resolveClass(solution, "KthLargest", null);
    if (!KthLargestClass) {
        throw new Error("missing KthLargest class");
    }
    const inp = test.input || {};
    const obj = new KthLargestClass(inp.k, inp.nums || []);
    const operations = inp.operations || [];
    const values = inp.values || [];
    const out = [];
    for (let i = 0; i < operations.length; i += 1) {
        const op = operations[i];
        const val = values[i];
        if (op !== "add") {
            throw new Error(`unknown op ${op}`);
        }
        out.push(obj.add(val));
    }
    gotValue = out;
"""
    elif kind == "class_min_stack":
        call_block = """
    const MinStackClass = resolveClass(solution, "MinStack", null);
    if (!MinStackClass) {
        throw new Error("missing MinStack class");
    }
    const inp = test.input || {};
    const operations = inp.operations || [];
    const values = inp.values || [];
    const obj = new MinStackClass();
    const out = [];
    for (let i = 0; i < operations.length; i += 1) {
        const op = operations[i];
        const val = values[i];
        if (op === "push") {
            obj.push(val);
            out.push(null);
        } else if (op === "pop") {
            obj.pop();
            out.push(null);
        } else if (op === "top") {
            out.push(obj.top());
        } else if (op === "get_min") {
            out.push(obj.getMin());
        } else {
            throw new Error(`unknown op ${op}`);
        }
    }
    gotValue = out;
"""
    elif kind == "class_trie":
        call_block = """
    const TrieClass = resolveClass(solution, "Trie", null);
    if (!TrieClass) {
        throw new Error("missing Trie class");
    }
    const inp = test.input || {};
    const operations = inp.operations || [];
    const values = inp.values || [];
    const obj = new TrieClass();
    const out = [];
    for (let i = 0; i < operations.length; i += 1) {
        const op = operations[i];
        const val = values[i];
        if (op === "insert") {
            obj.insert(val);
            out.push(null);
        } else if (op === "search") {
            out.push(obj.search(val));
        } else if (op === "starts_with") {
            out.push(obj.startsWith(val));
        } else {
            throw new Error(`unknown op ${op}`);
        }
    }
    gotValue = out;
"""
    elif kind == "linked_list_cycle":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const inp = test.input || {};
    const values = inp.list || [];
    const pos = inp.pos == null ? -1 : inp.pos;
    const head = arrayToLinkedList(values, NodeClass);
    if (pos >= 0) {
        let tail = head;
        let target = null;
        let idx = 0;
        while (tail && tail.next) {
            if (idx === pos) {
                target = tail;
            }
            tail = tail.next;
            idx += 1;
        }
        if (tail) {
            if (idx === pos) {
                target = tail;
            }
            if (target) {
                tail.next = target;
            }
        }
    }
    gotValue = fn(head);
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "linked_list_merge_k":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const raw = test.input || [];
    const lists = raw.map((arr) => arrayToLinkedList(arr, NodeClass));
    gotValue = linkedListToArray(fn(lists));
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "linked_list_merge_two":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const raw = test.input || [];
    const left = arrayToLinkedList(raw.length > 0 ? raw[0] : [], NodeClass);
    const right = arrayToLinkedList(raw.length > 1 ? raw[1] : [], NodeClass);
    gotValue = linkedListToArray(fn(left, right));
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "linked_list_remove_nth":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const inp = test.input || {};
    const head = arrayToLinkedList(inp.list || [], NodeClass);
    gotValue = linkedListToArray(fn(head, inp.n));
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "linked_list_reorder":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const head = arrayToLinkedList(test.input || [], NodeClass);
    const out = fn(head);
    gotValue = linkedListToArray(out == null ? head : out);
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "linked_list_reverse":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "ListNode", ListNodeFallback);
    const head = arrayToLinkedList(test.input || [], NodeClass);
    gotValue = linkedListToArray(fn(head));
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "tree_lca":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const TreeNodeClass = resolveClass(solution, "TreeNode", TreeNodeFallback);
    const inp = test.input || {};
    const root = arrayToTree(inp.root || [], TreeNodeClass);
    gotValue = fn(root, inp.p, inp.q);
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "clone_graph":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const NodeClass = resolveClass(solution, "Node", GraphNodeFallback);
    const inp = test.input || {};
    const root = adjListToGraph(inp.adjList || [], NodeClass);
    gotValue = graphToAdjList(fn(root));
""".replace("__FUNCTION_NAMES__", json.dumps(function_names))
    elif kind == "tree_plain":
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const TreeNodeClass = resolveClass(solution, "TreeNode", TreeNodeFallback);
    const inp = test.input || {};
    const args = [];
    for (const key of __ARG_KEYS__) {
        let value = inp[key];
        if (key === "root" || key === "subRoot") {
            value = arrayToTree(value, TreeNodeClass);
        } else if ((key === "p" || key === "q") && Array.isArray(value)) {
            value = arrayToTree(value, TreeNodeClass);
        }
        args.push(value);
    }
    const out = fn(...args);
    if (__CHALLENGE__ === "trees/build-tree-preorder" || __CHALLENGE__ === "trees/invert-binary-tree") {
        gotValue = treeToArray(out);
    } else {
        gotValue = out;
    }
""".replace("__FUNCTION_NAMES__", json.dumps(function_names)).replace("__ARG_KEYS__", json.dumps(arg_keys)).replace("__CHALLENGE__", json.dumps(challenge))
    else:
        call_block = """
    const fn = resolveFunction(solution, __FUNCTION_NAMES__);
    const caseInput = test.input;
    let args = [];
    if (isObject(caseInput)) {
        args = __ARG_KEYS__.map((key) => caseInput[key]);
    } else if (Array.isArray(caseInput) && __ARG_KEYS__.length > 1) {
        args.push(caseInput);
        for (let i = 1; i < __ARG_KEYS__.length; i += 1) {
            const key = __ARG_KEYS__[i];
            if (Object.prototype.hasOwnProperty.call(test, key)) {
                args.push(test[key]);
            } else {
                args.push(i < caseInput.length ? caseInput[i] : null);
            }
        }
    } else if (caseInput !== undefined) {
        args = [caseInput];
    } else {
        args = __ARG_KEYS__.map((key) => test[key]);
    }
    gotValue = fn(...args);
""".replace("__FUNCTION_NAMES__", json.dumps(function_names)).replace("__ARG_KEYS__", json.dumps(arg_keys))

    return f"""const fs = require(\"fs\");
const path = require(\"path\");

const TESTS = JSON.parse({json.dumps(tests_json)});
const COMPARE_MODE = {json.dumps(mode)};

function isObject(v) {{
    return v !== null && typeof v === \"object\" && !Array.isArray(v);
}}

function resolveFunction(solution, names) {{
    for (const name of names) {{
        if (typeof solution[name] === \"function\") {{
            return solution[name];
        }}
    }}
    throw new Error(\"missing function\");
}}

function resolveClass(solution, name, fallback) {{
    if (typeof solution[name] === \"function\") {{
        return solution[name];
    }}
    return fallback;
}}

class ListNodeFallback {{
    constructor(val = 0, next = null) {{
        this.val = val;
        this.next = next;
    }}
}}

class TreeNodeFallback {{
    constructor(val = 0, left = null, right = null) {{
        this.val = val;
        this.left = left;
        this.right = right;
    }}
}}

class GraphNodeFallback {{
    constructor(val = 0, neighbors = []) {{
        this.val = val;
        this.neighbors = neighbors;
    }}
}}

function arrayToLinkedList(values, NodeClass) {{
    let head = null;
    let tail = null;
    for (const value of values) {{
        const node = new NodeClass(value);
        if (head === null) {{
            head = node;
            tail = node;
        }} else {{
            tail.next = node;
            tail = node;
        }}
    }}
    return head;
}}

function linkedListToArray(head) {{
    const out = [];
    let node = head;
    let guard = 0;
    while (node !== null && node !== undefined && guard < 100000) {{
        out.push(node.val);
        node = node.next;
        guard += 1;
    }}
    return out;
}}

function arrayToTree(values, TreeNodeClass) {{
    if (!Array.isArray(values) || values.length === 0 || values[0] === null) {{
        return null;
    }}
    const nodes = values.map((v) => (v === null ? null : new TreeNodeClass(v)));
    let idx = 1;
    for (const node of nodes) {{
        if (node === null) {{
            continue;
        }}
        if (idx < nodes.length) {{
            node.left = nodes[idx];
            idx += 1;
        }}
        if (idx < nodes.length) {{
            node.right = nodes[idx];
            idx += 1;
        }}
    }}
    return nodes[0];
}}

function treeToArray(root) {{
    if (root === null || root === undefined) {{
        return [];
    }}
    const out = [];
    const queue = [root];
    let i = 0;
    while (i < queue.length) {{
        const node = queue[i];
        i += 1;
        if (node === null) {{
            out.push(null);
            continue;
        }}
        out.push(node.val);
        queue.push(node.left === undefined ? null : node.left);
        queue.push(node.right === undefined ? null : node.right);
    }}
    while (out.length > 0 && out[out.length - 1] === null) {{
        out.pop();
    }}
    return out;
}}

function adjListToGraph(adjList, NodeClass) {{
    if (!Array.isArray(adjList) || adjList.length === 0) {{
        return null;
    }}
    const nodes = [];
    for (let i = 0; i < adjList.length; i += 1) {{
        nodes.push(new NodeClass(i + 1, []));
    }}
    for (let i = 0; i < adjList.length; i += 1) {{
        nodes[i].neighbors = (adjList[i] || []).map((n) => nodes[n - 1]);
    }}
    return nodes[0];
}}

function graphToAdjList(root) {{
    if (root === null || root === undefined) {{
        return [];
    }}
    const seen = new Map();
    const queue = [root];
    while (queue.length > 0) {{
        const node = queue.shift();
        if (seen.has(node.val)) {{
            continue;
        }}
        seen.set(node.val, node);
        for (const nei of node.neighbors || []) {{
            queue.push(nei);
        }}
    }}
    let maxId = 0;
    for (const id of seen.keys()) {{
        if (id > maxId) {{
            maxId = id;
        }}
    }}
    const out = Array.from({{ length: maxId }}, () => []);
    for (const [id, node] of seen.entries()) {{
        const neighbors = (node.neighbors || []).map((n) => n.val).sort((a, b) => a - b);
        out[id - 1] = neighbors;
    }}
    return out;
}}

function canonical(value, mode) {{
    if (mode === \"pair_unordered\" || mode === \"list_unordered\" || mode === \"strings_unordered\") {{
        if (Array.isArray(value)) {{
            return [...value].sort();
        }}
        return value;
    }}
    if (mode === \"groups_unordered\") {{
        if (!Array.isArray(value)) {{
            return value;
        }}
        return value.map((group) => (Array.isArray(group) ? [...group].sort() : group)).sort((a, b) => JSON.stringify(a).localeCompare(JSON.stringify(b)));
    }}
    if (mode === \"nested_unordered\") {{
        if (!Array.isArray(value)) {{
            return value;
        }}
        return value.map((item) => (Array.isArray(item) ? [...item].sort((a, b) => (a > b ? 1 : a < b ? -1 : 0)) : item)).sort((a, b) => JSON.stringify(a).localeCompare(JSON.stringify(b)));
    }}
    return value;
}}

function equalValues(got, expected, mode) {{
    if (mode === \"pair_target_1idx\") {{
        return false;
    }}
    if (mode === \"float_sequence\") {{
        if (!Array.isArray(got) || !Array.isArray(expected) || got.length !== expected.length) {{
            return false;
        }}
        for (let i = 0; i < got.length; i += 1) {{
            const a = got[i];
            const b = expected[i];
            if (a === null && b === null) {{
                continue;
            }}
            if (a === null || b === null) {{
                return false;
            }}
            if (Math.abs(Number(a) - Number(b)) > 1e-9) {{
                return false;
            }}
        }}
        return true;
    }}
    return JSON.stringify(canonical(got, mode)) === JSON.stringify(canonical(expected, mode));
}}

function render(value) {{
    return JSON.stringify(value);
}}

function loadSolution(solutionPath) {{
    const resolved = path.resolve(solutionPath);
    delete require.cache[resolved];
    return require(resolved);
}}

function run(solutionPath, testsPath, challengePath) {{
    fs.readFileSync(testsPath, \"utf-8\");
    fs.readFileSync(challengePath, \"utf-8\");
    const solution = loadSolution(solutionPath);
    let failures = 0;
    for (let i = 0; i < TESTS.length; i += 1) {{
        const test = TESTS[i];
        let gotValue;
{call_block}
        const expected = test.expected;
        const expectedLength = test.expected_length;
        let passed = false;
        if (COMPARE_MODE === \"pair_target_1idx\") {{
            const inp = test.input || {{}};
            const numbers = Array.isArray(inp.numbers) ? inp.numbers : [];
            const target = inp.target;
            if (Array.isArray(gotValue) && gotValue.length === 2 && target !== undefined && target !== null) {{
                const i = Number(gotValue[0]) - 1;
                const j = Number(gotValue[1]) - 1;
                passed = i >= 0 && j >= 0 && i < numbers.length && j < numbers.length && i !== j && Number(numbers[i]) + Number(numbers[j]) === Number(target);
            }} else {{
                passed = false;
            }}
        }} else if (expectedLength !== undefined) {{
            passed = Array.isArray(gotValue) && gotValue.length === expectedLength;
        }} else {{
            passed = equalValues(gotValue, expected, COMPARE_MODE);
        }}
        const idx = i + 1;
        if (passed) {{
            console.log(`PASS ${{idx}}`);
        }} else {{
            failures += 1;
            if (expectedLength !== undefined) {{
                console.log(`FAIL ${{idx}} got=${{render(gotValue)}} expected=${{render(expectedLength)}}`);
            }} else {{
                console.log(`FAIL ${{idx}} got=${{render(gotValue)}} expected=${{render(expected)}}`);
            }}
        }}
    }}
    return failures === 0;
}}

function main() {{
    const builderDir = __dirname;
    const solutionPath = process.argv[2] ? path.resolve(process.argv[2]) : path.resolve(builderDir, \"..\", \"setup\", \"javascript.js\");
    const testsPath = process.argv[3] ? path.resolve(process.argv[3]) : path.resolve(builderDir, \"..\", \"tests.yaml\");
    const challengePath = process.argv[4] ? path.resolve(process.argv[4]) : path.resolve(builderDir, \"..\", \"challenge.yaml\");
    const ok = run(solutionPath, testsPath, challengePath);
    process.exit(ok ? 0 : 1);
}}

main();
"""


def find_name_case_insensitive(names, candidates):
    lowered = {n.lower(): n for n in names}
    for c in candidates:
        if c.lower() in lowered:
            return lowered[c.lower()]
    return None


def is_go_graph_type(ptype: str) -> bool:
    t = ptype.replace(" ", "")
    return t in ("*Node", "[]*Node")


def normalize_cpp_type(ptype: str):
    t = ptype.replace("const", "").replace("&", "").strip()
    t = re.sub(r"\s+", " ", t)
    t = t.replace("std:: ", "std::")
    t = t.replace(" *", "*").replace("* ", "*")
    return t


def cpp_convert_expr(ptype: str, raw_expr: str):
    t = normalize_cpp_type(ptype)
    mapping = {
        "int": f"toInt({raw_expr})",
        "double": f"toDoubleValue({raw_expr})",
        "bool": f"toBoolValue({raw_expr})",
        "std::string": f"toStringValue({raw_expr})",
        "string": f"toStringValue({raw_expr})",
        "std::vector<int>": f"toIntVector({raw_expr})",
        "std::vector<std::vector<int>>": f"toIntMatrix({raw_expr})",
        "std::vector<std::string>": f"toStringVector({raw_expr})",
        "std::vector<std::vector<std::string>>": f"toStringMatrix({raw_expr})",
        "std::vector<char>": f"toCharVector({raw_expr})",
        "std::vector<std::vector<char>>": f"toCharMatrix({raw_expr})",
        "ListNode*": f"buildList({raw_expr})",
        "std::vector<ListNode*>": f"buildListVector({raw_expr})",
        "TreeNode*": f"buildTree({raw_expr})",
        "Node*": f"buildGraph({raw_expr})",
    }
    if t not in mapping:
        raise RuntimeError(f"unsupported cpp type {ptype}")
    return mapping[t]


def cpp_result_expr(ret_type: str, result_var: str):
    t = normalize_cpp_type(ret_type)
    if t == "ListNode*":
        return f"listToValue({result_var})"
    if t == "TreeNode*":
        return f"treeToValue({result_var})"
    if t == "Node*":
        return f"graphToValue({result_var})"
    if t == "void":
        return "Value()"
    return f"toValue({result_var})"


def cpp_value_expr(v):
    if v is None:
        return "Value()"
    if isinstance(v, bool):
        return "Value(true)" if v else "Value(false)"
    if isinstance(v, (int, float)):
        return f"Value({float(v)!r})"
    if isinstance(v, str):
        return f"Value({json.dumps(v)})"
    if isinstance(v, list):
        inner = ", ".join(cpp_value_expr(item) for item in v)
        return f"Value::array({{{inner}}})"
    if isinstance(v, dict):
        inner = ", ".join(f"{{{json.dumps(str(k))}, {cpp_value_expr(val)}}}" for k, val in v.items())
        return f"Value::object({{{inner}}})"
    raise RuntimeError(f"unsupported cpp value {type(v)}")


def normalize_rust_type(ptype: str):
    t = re.sub(r"\s+", "", ptype)
    return t


def rust_convert_expr(ptype: str, raw_expr: str):
    t = normalize_rust_type(ptype)
    mapping = {
        "i32": f"to_i32(&{raw_expr})",
        "f64": f"to_f64(&{raw_expr})",
        "bool": f"to_bool(&{raw_expr})",
        "String": f"to_string_value(&{raw_expr})",
        "Vec<i32>": f"to_i32_vec(&{raw_expr})",
        "Vec<Vec<i32>>": f"to_i32_matrix(&{raw_expr})",
        "Vec<String>": f"to_string_vec(&{raw_expr})",
        "Vec<Vec<String>>": f"to_string_matrix(&{raw_expr})",
        "Vec<char>": f"to_char_vec(&{raw_expr})",
        "Vec<Vec<char>>": f"to_char_matrix(&{raw_expr})",
        "Option<Box<ListNode>>": f"build_list(&{raw_expr})",
        "Vec<Option<Box<ListNode>>>": f"build_list_vec(&{raw_expr})",
        "Option<Rc<RefCell<TreeNode>>>": f"build_tree(&{raw_expr})",
        "HashMap<i32,Vec<i32>>": f"to_graph_map(&{raw_expr})",
        "&HashMap<i32,Vec<i32>>": f"to_graph_map(&{raw_expr})",
    }
    if t not in mapping:
        raise RuntimeError(f"unsupported rust type {ptype}")
    return mapping[t]


def rust_result_expr(ret_type: str, result_var: str):
    t = normalize_rust_type(ret_type)
    if t in ("", "()"):
        return "Value::Null"
    if t == "Option<Box<ListNode>>":
        return f"list_to_value(&{result_var})"
    if t == "Option<Rc<RefCell<TreeNode>>>":
        return f"tree_to_value({result_var})"
    if t in ("HashMap<i32,Vec<i32>>",):
        return f"graph_map_to_value(&{result_var})"
    return f"to_value({result_var})"


def rust_value_expr(v):
    if v is None:
        return "Value::Null"
    if isinstance(v, bool):
        return "Value::Bool(true)" if v else "Value::Bool(false)"
    if isinstance(v, (int, float)):
        return f"Value::Num({float(v)!r})"
    if isinstance(v, str):
        return f"Value::Str({json.dumps(v)}.to_string())"
    if isinstance(v, list):
        inner = ", ".join(rust_value_expr(item) for item in v)
        return f"Value::Arr(vec![{inner}])"
    if isinstance(v, dict):
        inner = ", ".join(f"({json.dumps(str(k))}.to_string(), {rust_value_expr(val)})" for k, val in v.items())
        return f"obj(vec![{inner}])"
    raise RuntimeError(f"unsupported rust value {type(v)}")


def go_convert_expr(ptype, raw_expr):
    mapping = {
        "int": f"toInt({raw_expr})",
        "float64": f"toFloat({raw_expr})",
        "string": f"toString({raw_expr})",
        "bool": f"toBool({raw_expr})",
        "[]int": f"toIntSlice({raw_expr})",
        "[][]int": f"toIntMatrix({raw_expr})",
        "[]string": f"toStringSlice({raw_expr})",
        "[][]string": f"toStringMatrix({raw_expr})",
        "[]byte": f"toByteSlice({raw_expr})",
        "[][]byte": f"toByteMatrix({raw_expr})",
        "*ListNode": f"buildList({raw_expr})",
        "[]*ListNode": f"buildListArray({raw_expr})",
        "*TreeNode": f"buildTree({raw_expr})",
        "*Node": f"buildGraph({raw_expr})",
    }
    if ptype not in mapping:
        raise RuntimeError(f"unsupported go type {ptype}")
    return mapping[ptype]


def go_result_expr(ret_type, result_var):
    if ret_type == "*ListNode":
        return f"normalizeValue(listToAny({result_var}))"
    if ret_type == "*TreeNode":
        return f"normalizeValue(treeToAny({result_var}))"
    if ret_type == "*Node":
        return f"normalizeValue(graphToAny({result_var}))"
    return f"normalizeValue({result_var})"


def go_builder(challenge, tests, funcs, methods, arg_keys):
    kind = python_kind(challenge)
    mode = compare_mode(challenge)
    tests_json = json.dumps(tests, separators=(",", ":"), ensure_ascii=False)
    arg_keys_literal = ", ".join(json.dumps(k) for k in arg_keys)
    func_names = list(funcs.keys())
    method_names = list(methods.keys())

    need_list = kind in {
        "linked_list_cycle",
        "linked_list_merge_k",
        "linked_list_merge_two",
        "linked_list_remove_nth",
        "linked_list_reorder",
        "linked_list_reverse",
    }
    need_tree = kind in {"tree_plain", "tree_lca", "serialize_deserialize"}
    need_graph = kind == "clone_graph"

    call_block = ""

    if kind == "encode_decode":
        encode_name = find_name_case_insensitive(func_names, ["Encode"])
        decode_name = find_name_case_insensitive(func_names, ["Decode"])
        call_block = f"""
        inp := getInputMap(test)
        raw := inp["strs"]
        gotValue = normalizeValue({decode_name}({encode_name}(toStringSlice(raw))))
"""
    elif kind == "serialize_deserialize":
        serialize_name = find_name_case_insensitive(func_names, ["Serialize"])
        deserialize_name = find_name_case_insensitive(func_names, ["Deserialize"])
        call_block = f"""
        inp := getInputMap(test)
        root := buildTree(inp["root"])
        gotValue = normalizeValue(treeToAny({deserialize_name}({serialize_name}(root))))
"""
    elif kind == "class_min_stack":
        ctor = find_name_case_insensitive(func_names, ["NewMinStack"])
        push = find_name_case_insensitive(method_names, ["Push"])
        pop = find_name_case_insensitive(method_names, ["Pop"])
        top = find_name_case_insensitive(method_names, ["Top"])
        getmin = find_name_case_insensitive(method_names, ["GetMin"])
        call_block = f"""
        inp := getInputMap(test)
        operations := toStringSlice(inp["operations"])
        values := toAnySlice(inp["values"])
        obj := {ctor}()
        out := make([]any, 0, len(operations))
        for i, op := range operations {{
            var val any
            if i < len(values) {{
                val = values[i]
            }}
            switch op {{
            case "push":
                obj.{push}(toInt(val))
                out = append(out, nil)
            case "pop":
                obj.{pop}()
                out = append(out, nil)
            case "top":
                out = append(out, obj.{top}())
            case "get_min":
                out = append(out, obj.{getmin}())
            default:
                panic("unknown op")
            }}
        }}
        gotValue = normalizeValue(out)
"""
    elif kind == "class_kth_largest":
        ctor = find_name_case_insensitive(func_names, ["NewKthLargest"])
        addm = find_name_case_insensitive(method_names, ["Add"])
        call_block = f"""
        inp := getInputMap(test)
        obj := {ctor}(toInt(inp["k"]), toIntSlice(inp["nums"]))
        operations := toStringSlice(inp["operations"])
        values := toAnySlice(inp["values"])
        out := make([]any, 0, len(operations))
        for i, op := range operations {{
            if op != "add" {{
                panic("unknown op")
            }}
            var val any
            if i < len(values) {{
                val = values[i]
            }}
            out = append(out, obj.{addm}(toInt(val)))
        }}
        gotValue = normalizeValue(out)
"""
    elif kind == "class_median_finder":
        ctor = find_name_case_insensitive(func_names, ["NewMedianFinder"])
        addm = find_name_case_insensitive(method_names, ["AddNum"])
        findm = find_name_case_insensitive(method_names, ["FindMedian"])
        call_block = f"""
        inp := getInputMap(test)
        obj := {ctor}()
        operations := toStringSlice(inp["operations"])
        values := toAnySlice(inp["values"])
        out := make([]any, 0, len(operations))
        for i, op := range operations {{
            var val any
            if i < len(values) {{
                val = values[i]
            }}
            switch op {{
            case "add_num":
                obj.{addm}(toInt(val))
                out = append(out, nil)
            case "find_median":
                out = append(out, obj.{findm}())
            default:
                panic("unknown op")
            }}
        }}
        gotValue = normalizeValue(out)
"""
    elif kind == "class_trie":
        ctor = find_name_case_insensitive(func_names, ["NewTrie"])
        insert = find_name_case_insensitive(method_names, ["Insert"])
        search = find_name_case_insensitive(method_names, ["Search"])
        starts = find_name_case_insensitive(method_names, ["StartsWith"])
        call_block = f"""
        inp := getInputMap(test)
        obj := {ctor}()
        operations := toStringSlice(inp["operations"])
        values := toAnySlice(inp["values"])
        out := make([]any, 0, len(operations))
        for i, op := range operations {{
            val := ""
            if i < len(values) {{
                val = toString(values[i])
            }}
            switch op {{
            case "insert":
                obj.{insert}(val)
                out = append(out, nil)
            case "search":
                out = append(out, obj.{search}(val))
            case "starts_with":
                out = append(out, obj.{starts}(val))
            default:
                panic("unknown op")
            }}
        }}
        gotValue = normalizeValue(out)
"""
    elif kind == "linked_list_cycle":
        fname = func_names[0]
        call_block = f"""
        inp := getInputMap(test)
        values := toIntSlice(inp["list"])
        pos := toInt(inp["pos"])
        head := buildCycleList(values, pos)
        gotValue = normalizeValue({fname}(head))
"""
    elif kind == "linked_list_merge_k":
        fname = func_names[0]
        call_block = f"""
        gotValue = normalizeValue(listToAny({fname}(buildListArray(test["input"]))))
"""
    elif kind == "linked_list_merge_two":
        fname = func_names[0]
        call_block = f"""
        raw := toAnySlice(test["input"])
        var left any = []any{{}}
        var right any = []any{{}}
        if len(raw) > 0 {{
            left = raw[0]
        }}
        if len(raw) > 1 {{
            right = raw[1]
        }}
        gotValue = normalizeValue(listToAny({fname}(buildList(left), buildList(right))))
"""
    elif kind == "linked_list_remove_nth":
        fname = func_names[0]
        call_block = f"""
        inp := getInputMap(test)
        gotValue = normalizeValue(listToAny({fname}(buildList(inp["list"]), toInt(inp["n"]))))
"""
    elif kind == "linked_list_reorder":
        fname = func_names[0]
        call_block = f"""
        head := buildList(test["input"])
        {fname}(head)
        gotValue = normalizeValue(listToAny(head))
"""
    elif kind == "linked_list_reverse":
        fname = func_names[0]
        call_block = f"""
        gotValue = normalizeValue(listToAny({fname}(buildList(test["input"]))))
"""
    elif kind == "tree_lca":
        fname = func_names[0]
        call_block = f"""
        inp := getInputMap(test)
        gotValue = normalizeValue({fname}(buildTree(inp["root"]), toInt(inp["p"]), toInt(inp["q"])))
"""
    elif kind == "clone_graph":
        fname = func_names[0]
        call_block = f"""
        inp := getInputMap(test)
        gotValue = normalizeValue(graphToAny({fname}(buildGraph(inp["adjList"]))))
"""
    else:
        if not func_names:
            raise RuntimeError(f"no go function for {challenge}")
        fname = func_names[0]
        sig = funcs[fname]
        params = sig["params"]
        ret = sig["ret"]
        lines = []
        for idx, (_, ptype) in enumerate(params):
            lines.append(f"        raw{idx} := getArg(test, argKeys, {idx})")
            lines.append(f"        arg{idx} := {go_convert_expr(ptype, f'raw{idx}')}")
        args_join = ", ".join([f"arg{i}" for i in range(len(params))])
        lines.append(f"        result := {fname}({args_join})")
        lines.append(f"        gotValue = {go_result_expr(ret, 'result')}")
        call_block = "\n".join(lines) + "\n"
        if "ListNode" in ret or any("ListNode" in p for _, p in params):
            need_list = True
        if "TreeNode" in ret or any("TreeNode" in p for _, p in params):
            need_tree = True
        if is_go_graph_type(ret) or any(is_go_graph_type(p) for _, p in params):
            need_graph = True

    list_helpers = ""
    if need_list:
        list_helpers = """
func buildList(v any) *ListNode {
    values := toIntSlice(v)
    var head *ListNode
    var tail *ListNode
    for _, value := range values {
        node := &ListNode{Val: value}
        if head == nil {
            head = node
            tail = node
        } else {
            tail.Next = node
            tail = node
        }
    }
    return head
}

func buildListArray(v any) []*ListNode {
    raw := toAnySlice(v)
    out := make([]*ListNode, 0, len(raw))
    for _, item := range raw {
        out = append(out, buildList(item))
    }
    return out
}

func listToAny(head *ListNode) []any {
    out := make([]any, 0)
    cur := head
    guard := 0
    for cur != nil && guard < 100000 {
        out = append(out, cur.Val)
        cur = cur.Next
        guard++
    }
    return out
}

func buildCycleList(values []int, pos int) *ListNode {
    var head *ListNode
    var tail *ListNode
    var target *ListNode
    for i, value := range values {
        node := &ListNode{Val: value}
        if head == nil {
            head = node
            tail = node
        } else {
            tail.Next = node
            tail = node
        }
        if i == pos {
            target = node
        }
    }
    if tail != nil && target != nil && pos >= 0 {
        tail.Next = target
    }
    return head
}
"""

    tree_helpers = ""
    if need_tree:
        tree_helpers = """
func buildTree(v any) *TreeNode {
    raw := toAnySlice(v)
    if len(raw) == 0 || raw[0] == nil {
        return nil
    }
    nodes := make([]*TreeNode, len(raw))
    for i, item := range raw {
        if item == nil {
            continue
        }
        nodes[i] = &TreeNode{Val: toInt(item)}
    }
    idx := 1
    for _, node := range nodes {
        if node == nil {
            continue
        }
        if idx < len(nodes) {
            node.Left = nodes[idx]
            idx++
        }
        if idx < len(nodes) {
            node.Right = nodes[idx]
            idx++
        }
    }
    return nodes[0]
}

func treeToAny(root *TreeNode) []any {
    if root == nil {
        return []any{}
    }
    queue := []*TreeNode{root}
    out := make([]any, 0)
    for i := 0; i < len(queue); i++ {
        node := queue[i]
        if node == nil {
            out = append(out, nil)
            continue
        }
        out = append(out, node.Val)
        queue = append(queue, node.Left, node.Right)
    }
    for len(out) > 0 && out[len(out)-1] == nil {
        out = out[:len(out)-1]
    }
    return out
}
"""

    graph_helpers = ""
    if need_graph:
        graph_helpers = """
func buildGraph(v any) *Node {
    adj := toIntMatrix(v)
    if len(adj) == 0 {
        return nil
    }
    nodes := make([]*Node, len(adj))
    for i := range adj {
        nodes[i] = &Node{Val: i + 1}
    }
    for i, neighbors := range adj {
        items := make([]*Node, 0, len(neighbors))
        for _, n := range neighbors {
            if n >= 1 && n <= len(nodes) {
                items = append(items, nodes[n-1])
            }
        }
        nodes[i].Neighbors = items
    }
    return nodes[0]
}

func graphToAny(root *Node) []any {
    if root == nil {
        return []any{}
    }
    seen := map[int]*Node{}
    queue := []*Node{root}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        if node == nil {
            continue
        }
        if _, ok := seen[node.Val]; ok {
            continue
        }
        seen[node.Val] = node
        for _, nei := range node.Neighbors {
            queue = append(queue, nei)
        }
    }
    maxID := 0
    for id := range seen {
        if id > maxID {
            maxID = id
        }
    }
    out := make([]any, maxID)
    for id, node := range seen {
        row := make([]int, 0, len(node.Neighbors))
        for _, nei := range node.Neighbors {
            row = append(row, nei.Val)
        }
        sort.Ints(row)
        anyRow := make([]any, 0, len(row))
        for _, value := range row {
            anyRow = append(anyRow, value)
        }
        out[id-1] = anyRow
    }
    for i := range out {
        if out[i] == nil {
            out[i] = []any{}
        }
    }
    return out
}
"""

    return f"""package main

import (
    "encoding/json"
    "fmt"
    "math"
    "os"
    "reflect"
    "sort"
)

const testsJSON = {json.dumps(tests_json)}
const compareMode = {json.dumps(mode)}

func main() {{
    testsPath := "../tests.yaml"
    challengePath := "../challenge.yaml"
    if len(os.Args) > 3 {{
        testsPath = os.Args[3]
    }}
    if len(os.Args) > 4 {{
        challengePath = os.Args[4]
    }}
    _, _ = os.ReadFile(testsPath)
    _, _ = os.ReadFile(challengePath)

    var tests []map[string]any
    if err := json.Unmarshal([]byte(testsJSON), &tests); err != nil {{
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(1)
    }}

    argKeys := []string{{{arg_keys_literal}}}
    _ = argKeys
    failed := false
    for i, test := range tests {{
        var gotValue any
{call_block}
        expected := test["expected"]
        expectedLength, hasExpectedLength := test["expected_length"]
        passed := false
        if hasExpectedLength {{
            if arr, ok := gotValue.([]any); ok {{
                passed = len(arr) == toInt(expectedLength)
            }} else {{
                passed = false
            }}
        }} else if compareMode == "pair_target_1idx" {{
            inputObj := getInputMap(test)
            numbers := toAnySlice(inputObj["numbers"])
            target := toInt(inputObj["target"])
            if pair, ok := gotValue.([]any); ok && len(pair) == 2 {{
                i1 := toInt(pair[0]) - 1
                i2 := toInt(pair[1]) - 1
                if i1 >= 0 && i2 >= 0 && i1 < len(numbers) && i2 < len(numbers) && i1 != i2 {{
                    passed = toInt(numbers[i1])+toInt(numbers[i2]) == target
                }} else {{
                    passed = false
                }}
            }} else {{
                passed = false
            }}
        }} else {{
            passed = equalValues(gotValue, expected, compareMode)
        }}
        idx := i + 1
        if passed {{
            fmt.Printf("PASS %d\\n", idx)
        }} else {{
            failed = true
            if hasExpectedLength {{
                fmt.Printf("FAIL %d got=%s expected=%s\\n", idx, render(gotValue), render(normalizeValue(expectedLength)))
            }} else {{
                fmt.Printf("FAIL %d got=%s expected=%s\\n", idx, render(gotValue), render(expected))
            }}
        }}
    }}
    if failed {{
        os.Exit(1)
    }}
}}

func getInputMap(test map[string]any) map[string]any {{
    if input, ok := test["input"]; ok {{
        if m, ok := input.(map[string]any); ok {{
            return m
        }}
    }}
    return map[string]any{{}}
}}

func getArg(test map[string]any, argKeys []string, idx int) any {{
    input, hasInput := test["input"]
    if hasInput {{
        if m, ok := input.(map[string]any); ok {{
            if idx < len(argKeys) {{
                return m[argKeys[idx]]
            }}
            return nil
        }}
        if arr, ok := input.([]any); ok {{
            if len(argKeys) > 1 {{
                if idx == 0 {{
                    return arr
                }}
                key := argKeys[idx]
                if v, ok := test[key]; ok {{
                    return v
                }}
                if idx < len(arr) {{
                    return arr[idx]
                }}
                return nil
            }}
            return arr
        }}
        return input
    }}
    if idx < len(argKeys) {{
        return test[argKeys[idx]]
    }}
    return nil
}}

func toAnySlice(v any) []any {{
    if v == nil {{
        return []any{{}}
    }}
    if arr, ok := v.([]any); ok {{
        return arr
    }}
    return []any{{}}
}}

func toInt(v any) int {{
    switch value := v.(type) {{
    case nil:
        return 0
    case int:
        return value
    case int64:
        return int(value)
    case float64:
        return int(value)
    case json.Number:
        n, _ := value.Int64()
        return int(n)
    default:
        return 0
    }}
}}

func toFloat(v any) float64 {{
    switch value := v.(type) {{
    case nil:
        return 0
    case float64:
        return value
    case int:
        return float64(value)
    case json.Number:
        n, _ := value.Float64()
        return n
    default:
        return 0
    }}
}}

func toString(v any) string {{
    switch value := v.(type) {{
    case nil:
        return ""
    case string:
        return value
    default:
        return fmt.Sprint(value)
    }}
}}

func toBool(v any) bool {{
    if b, ok := v.(bool); ok {{
        return b
    }}
    return false
}}

func toIntSlice(v any) []int {{
    raw := toAnySlice(v)
    out := make([]int, 0, len(raw))
    for _, item := range raw {{
        out = append(out, toInt(item))
    }}
    return out
}}

func toIntMatrix(v any) [][]int {{
    raw := toAnySlice(v)
    out := make([][]int, 0, len(raw))
    for _, item := range raw {{
        out = append(out, toIntSlice(item))
    }}
    return out
}}

func toStringSlice(v any) []string {{
    raw := toAnySlice(v)
    out := make([]string, 0, len(raw))
    for _, item := range raw {{
        out = append(out, toString(item))
    }}
    return out
}}

func toStringMatrix(v any) [][]string {{
    raw := toAnySlice(v)
    out := make([][]string, 0, len(raw))
    for _, item := range raw {{
        out = append(out, toStringSlice(item))
    }}
    return out
}}

func toByteSlice(v any) []byte {{
    raw := toAnySlice(v)
    out := make([]byte, 0, len(raw))
    for _, item := range raw {{
        s := toString(item)
        if s == "" {{
            out = append(out, 0)
        }} else {{
            out = append(out, s[0])
        }}
    }}
    return out
}}

func toByteMatrix(v any) [][]byte {{
    raw := toAnySlice(v)
    out := make([][]byte, 0, len(raw))
    for _, item := range raw {{
        out = append(out, toByteSlice(item))
    }}
    return out
}}

func normalizeValue(v any) any {{
    data, err := json.Marshal(v)
    if err != nil {{
        return v
    }}
    var out any
    if err := json.Unmarshal(data, &out); err != nil {{
        return v
    }}
    return out
}}

func render(v any) string {{
    normalized := normalizeValue(v)
    data, err := json.Marshal(normalized)
    if err != nil {{
        return "null"
    }}
    return string(data)
}}

func sortPrimitiveSlice(values []any) []any {{
    out := append([]any{{}}, values...)
    sort.Slice(out, func(i, j int) bool {{
        return fmt.Sprint(out[i]) < fmt.Sprint(out[j])
    }})
    return out
}}

func canonical(v any, mode string) any {{
    if mode == "pair_unordered" || mode == "list_unordered" || mode == "strings_unordered" {{
        if arr, ok := v.([]any); ok {{
            return sortPrimitiveSlice(arr)
        }}
        return v
    }}
    if mode == "groups_unordered" || mode == "nested_unordered" {{
        arr, ok := v.([]any)
        if !ok {{
            return v
        }}
        outer := make([]any, 0, len(arr))
        for _, item := range arr {{
            if inner, ok := item.([]any); ok {{
                outer = append(outer, sortPrimitiveSlice(inner))
            }} else {{
                outer = append(outer, item)
            }}
        }}
        sort.Slice(outer, func(i, j int) bool {{
            return render(outer[i]) < render(outer[j])
        }})
        return outer
    }}
    return v
}}

func equalFloatSequence(got, expected any) bool {{
    ga, ok1 := got.([]any)
    ea, ok2 := expected.([]any)
    if !ok1 || !ok2 || len(ga) != len(ea) {{
        return false
    }}
    for i := range ga {{
        if ga[i] == nil && ea[i] == nil {{
            continue
        }}
        if ga[i] == nil || ea[i] == nil {{
            return false
        }}
        if math.Abs(toFloat(ga[i])-toFloat(ea[i])) > 1e-9 {{
            return false
        }}
    }}
    return true
}}

func equalValues(got, expected any, mode string) bool {{
    g := normalizeValue(got)
    e := normalizeValue(expected)
    if mode == "pair_target_1idx" {{
        return false
    }}
    if mode == "float_sequence" {{
        return equalFloatSequence(g, e)
    }}
    cg := canonical(g, mode)
    ce := canonical(e, mode)
    if reflect.DeepEqual(cg, ce) {{
        return true
    }}
    if cg == nil {{
        if arr, ok := ce.([]any); ok && len(arr) == 0 {{
            return true
        }}
    }}
    if ce == nil {{
        if arr, ok := cg.([]any); ok && len(arr) == 0 {{
            return true
        }}
    }}
    return false
}}

{list_helpers}
{tree_helpers}
{graph_helpers}
"""


def cpp_builder(challenge, tests, solution_methods, class_methods, arg_keys):
    kind = python_kind(challenge)
    mode = compare_mode(challenge)
    tests_literal = ",\n        ".join(cpp_value_expr(t) for t in tests)
    arg_keys_literal = ", ".join(json.dumps(k) for k in arg_keys)

    need_list = kind in {
        "linked_list_cycle",
        "linked_list_merge_k",
        "linked_list_merge_two",
        "linked_list_remove_nth",
        "linked_list_reorder",
        "linked_list_reverse",
    }
    need_tree = kind in {"tree_plain", "tree_lca", "serialize_deserialize"}
    need_graph = kind == "clone_graph"

    call_block = ""

    if kind == "encode_decode":
        encode_name = find_name_case_insensitive(list(solution_methods.keys()), ["encode"])
        decode_name = find_name_case_insensitive(list(solution_methods.keys()), ["decode"])
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        gotValue = toValue(sol.{decode_name}(sol.{encode_name}(toStringVector(getField(input, "strs")))));
"""
    elif kind == "serialize_deserialize":
        serialize_name = find_name_case_insensitive(list(solution_methods.keys()), ["serialize"])
        deserialize_name = find_name_case_insensitive(list(solution_methods.keys()), ["deserialize"])
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        TreeNode* root = buildTree(getField(input, "root"));
        gotValue = treeToValue(sol.{deserialize_name}(sol.{serialize_name}(root)));
"""
    elif kind == "class_min_stack":
        methods_map = class_methods.get("MinStack", {})
        push = find_name_case_insensitive(list(methods_map.keys()), ["push"])
        pop = find_name_case_insensitive(list(methods_map.keys()), ["pop"])
        top = find_name_case_insensitive(list(methods_map.keys()), ["top"])
        get_min = find_name_case_insensitive(list(methods_map.keys()), ["getMin", "get_min"])
        call_block = f"""
        Value::Object input = asObject(getField(testObj, "input"));
        std::vector<std::string> operations = toStringVector(getField(input, "operations"));
        Value::Array values = asArray(getField(input, "values"));
        MinStack obj;
        Value::Array out;
        for (size_t j = 0; j < operations.size(); ++j) {{
            std::string op = operations[j];
            Value val = j < values.size() ? values[j] : Value();
            if (op == "push") {{
                obj.{push}(toInt(val));
                out.push_back(Value());
            }} else if (op == "pop") {{
                obj.{pop}();
                out.push_back(Value());
            }} else if (op == "top") {{
                out.push_back(toValue(obj.{top}()));
            }} else if (op == "get_min") {{
                out.push_back(toValue(obj.{get_min}()));
            }} else {{
                throw std::runtime_error("unknown op");
            }}
        }}
        gotValue = Value(out);
"""
    elif kind == "class_kth_largest":
        methods_map = class_methods.get("KthLargest", {})
        add_name = find_name_case_insensitive(list(methods_map.keys()), ["add"])
        call_block = f"""
        Value::Object input = asObject(getField(testObj, "input"));
        KthLargest obj(toInt(getField(input, "k")), toIntVector(getField(input, "nums")));
        std::vector<std::string> operations = toStringVector(getField(input, "operations"));
        Value::Array values = asArray(getField(input, "values"));
        Value::Array out;
        for (size_t j = 0; j < operations.size(); ++j) {{
            if (operations[j] != "add") {{
                throw std::runtime_error("unknown op");
            }}
            Value val = j < values.size() ? values[j] : Value();
            out.push_back(toValue(obj.{add_name}(toInt(val))));
        }}
        gotValue = Value(out);
"""
    elif kind == "class_median_finder":
        methods_map = class_methods.get("MedianFinder", {})
        add_name = find_name_case_insensitive(list(methods_map.keys()), ["addNum", "add_num"])
        find_name = find_name_case_insensitive(list(methods_map.keys()), ["findMedian", "find_median"])
        call_block = f"""
        Value::Object input = asObject(getField(testObj, "input"));
        MedianFinder obj;
        std::vector<std::string> operations = toStringVector(getField(input, "operations"));
        Value::Array values = asArray(getField(input, "values"));
        Value::Array out;
        for (size_t j = 0; j < operations.size(); ++j) {{
            std::string op = operations[j];
            Value val = j < values.size() ? values[j] : Value();
            if (op == "add_num") {{
                obj.{add_name}(toInt(val));
                out.push_back(Value());
            }} else if (op == "find_median") {{
                out.push_back(toValue(obj.{find_name}()));
            }} else {{
                throw std::runtime_error("unknown op");
            }}
        }}
        gotValue = Value(out);
"""
    elif kind == "class_trie":
        methods_map = class_methods.get("Trie", {})
        insert = find_name_case_insensitive(list(methods_map.keys()), ["insert"])
        search = find_name_case_insensitive(list(methods_map.keys()), ["search"])
        starts = find_name_case_insensitive(list(methods_map.keys()), ["startsWith", "starts_with"])
        call_block = f"""
        Value::Object input = asObject(getField(testObj, "input"));
        Trie obj;
        std::vector<std::string> operations = toStringVector(getField(input, "operations"));
        Value::Array values = asArray(getField(input, "values"));
        Value::Array out;
        for (size_t j = 0; j < operations.size(); ++j) {{
            std::string op = operations[j];
            std::string val = j < values.size() ? toStringValue(values[j]) : "";
            if (op == "insert") {{
                obj.{insert}(val);
                out.push_back(Value());
            }} else if (op == "search") {{
                out.push_back(toValue(obj.{search}(val)));
            }} else if (op == "starts_with") {{
                out.push_back(toValue(obj.{starts}(val)));
            }} else {{
                throw std::runtime_error("unknown op");
            }}
        }}
        gotValue = Value(out);
"""
    elif kind == "linked_list_cycle":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        ListNode* head = buildCycleList(toIntVector(getField(input, "list")), toInt(getField(input, "pos")));
        gotValue = toValue(sol.{fname}(head));
"""
    elif kind == "linked_list_merge_k":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        gotValue = listToValue(sol.{fname}(buildListVector(getField(testObj, "input"))));
"""
    elif kind == "linked_list_merge_two":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        Value::Array raw = asArray(getField(testObj, "input"));
        Value left = raw.size() > 0 ? raw[0] : Value::array({{}});
        Value right = raw.size() > 1 ? raw[1] : Value::array({{}});
        gotValue = listToValue(sol.{fname}(buildList(left), buildList(right)));
"""
    elif kind == "linked_list_remove_nth":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        gotValue = listToValue(sol.{fname}(buildList(getField(input, "list")), toInt(getField(input, "n"))));
"""
    elif kind == "linked_list_reorder":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        ListNode* head = buildList(getField(testObj, "input"));
        sol.{fname}(head);
        gotValue = listToValue(head);
"""
    elif kind == "linked_list_reverse":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        gotValue = listToValue(sol.{fname}(buildList(getField(testObj, "input"))));
"""
    elif kind == "tree_lca":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        gotValue = toValue(sol.{fname}(buildTree(getField(input, "root")), toInt(getField(input, "p")), toInt(getField(input, "q"))));
"""
    elif kind == "clone_graph":
        fname = next(iter(solution_methods.keys()))
        call_block = f"""
        Solution sol;
        Value::Object input = asObject(getField(testObj, "input"));
        gotValue = graphToValue(sol.{fname}(buildGraph(getField(input, "adjList"))));
"""
    else:
        if not solution_methods:
            raise RuntimeError(f"no cpp solution method for {challenge}")
        fname = next(iter(solution_methods.keys()))
        sig = solution_methods[fname]
        params = sig["params"]
        ret = sig["ret"]
        lines = ["        Solution sol;"]
        for idx, (_, ptype) in enumerate(params):
            lines.append(f"        Value raw{idx} = getArg(testObj, argKeys, {idx});")
            lines.append(f"        auto arg{idx} = {cpp_convert_expr(ptype, f'raw{idx}')};")
        args_join = ", ".join(f"arg{i}" for i in range(len(params)))
        nret = normalize_cpp_type(ret)
        if nret == "void":
            lines.append(f"        sol.{fname}({args_join});")
            lines.append("        gotValue = Value();")
        else:
            lines.append(f"        auto result = sol.{fname}({args_join});")
            lines.append(f"        gotValue = {cpp_result_expr(ret, 'result')};")
        call_block = "\n".join(lines) + "\n"
        if "ListNode" in nret or any("ListNode" in normalize_cpp_type(p) for _, p in params):
            need_list = True
        if "TreeNode" in nret or any("TreeNode" in normalize_cpp_type(p) for _, p in params):
            need_tree = True
        if "Node" in nret or any("Node" in normalize_cpp_type(p) for _, p in params):
            need_graph = True

    list_helpers = ""
    if need_list:
        list_helpers = """
ListNode* buildList(const Value& v) {
    std::vector<int> values = toIntVector(v);
    ListNode* head = nullptr;
    ListNode* tail = nullptr;
    for (int value : values) {
        ListNode* node = new ListNode(value);
        if (head == nullptr) {
            head = node;
            tail = node;
        } else {
            tail->next = node;
            tail = node;
        }
    }
    return head;
}

std::vector<ListNode*> buildListVector(const Value& v) {
    std::vector<ListNode*> out;
    for (const Value& item : asArray(v)) {
        out.push_back(buildList(item));
    }
    return out;
}

Value listToValue(ListNode* head) {
    Value::Array out;
    int guard = 0;
    ListNode* cur = head;
    while (cur != nullptr && guard < 100000) {
        out.push_back(toValue(cur->val));
        cur = cur->next;
        ++guard;
    }
    return Value(out);
}

ListNode* buildCycleList(const std::vector<int>& values, int pos) {
    ListNode* head = nullptr;
    ListNode* tail = nullptr;
    ListNode* target = nullptr;
    for (size_t i = 0; i < values.size(); ++i) {
        ListNode* node = new ListNode(values[i]);
        if (head == nullptr) {
            head = node;
            tail = node;
        } else {
            tail->next = node;
            tail = node;
        }
        if (static_cast<int>(i) == pos) {
            target = node;
        }
    }
    if (tail != nullptr && target != nullptr && pos >= 0) {
        tail->next = target;
    }
    return head;
}
"""

    tree_helpers = ""
    if need_tree:
        tree_helpers = """
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
"""

    graph_helpers = ""
    if need_graph:
        graph_helpers = """
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
"""

    return f"""#include <algorithm>
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

struct Value {{
    using Array = std::vector<Value>;
    using Object = std::map<std::string, Value>;
    std::variant<std::nullptr_t, bool, double, std::string, Array, Object> data;
    Value() : data(nullptr) {{}}
    Value(std::nullptr_t) : data(nullptr) {{}}
    Value(bool v) : data(v) {{}}
    Value(double v) : data(v) {{}}
    Value(const std::string& v) : data(v) {{}}
    Value(const char* v) : data(std::string(v)) {{}}
    Value(const Array& v) : data(v) {{}}
    Value(const Object& v) : data(v) {{}}
    bool isNull() const {{ return std::holds_alternative<std::nullptr_t>(data); }}
    static Value array(std::initializer_list<Value> init) {{ return Value(Array(init)); }}
    static Value object(std::initializer_list<std::pair<std::string, Value>> init) {{
        Object obj;
        for (const auto& item : init) {{
            obj[item.first] = item.second;
        }}
        return Value(obj);
    }}
}};

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

Value toValue(int v) {{ return Value(static_cast<double>(v)); }}
Value toValue(double v) {{ return Value(v); }}
Value toValue(bool v) {{ return Value(v); }}
Value toValue(const std::string& v) {{ return Value(v); }}
Value toValue(const char* v) {{ return Value(std::string(v)); }}
Value toValue(const Value& v) {{ return v; }}

template <typename T>
Value toValue(const std::vector<T>& values) {{
    Value::Array out;
    for (const auto& item : values) {{
        out.push_back(toValue(item));
    }}
    return Value(out);
}}

const Value::Array& asArray(const Value& v) {{
    static const Value::Array empty;
    if (auto ptr = std::get_if<Value::Array>(&v.data)) {{
        return *ptr;
    }}
    return empty;
}}

Value::Object asObject(const Value& v) {{
    if (auto ptr = std::get_if<Value::Object>(&v.data)) {{
        return *ptr;
    }}
    return Value::Object{{}};
}}

Value getField(const Value::Object& obj, const std::string& key) {{
    auto it = obj.find(key);
    if (it == obj.end()) {{
        return Value();
    }}
    return it->second;
}}

bool hasField(const Value::Object& obj, const std::string& key) {{
    return obj.find(key) != obj.end();
}}

std::string numberToString(double value) {{
    if (std::fabs(value - std::round(value)) < 1e-9) {{
        long long iv = static_cast<long long>(std::llround(value));
        return std::to_string(iv);
    }}
    std::ostringstream oss;
    oss << std::setprecision(15) << value;
    std::string out = oss.str();
    while (!out.empty() && out.back() == '0' && out.find('.') != std::string::npos) {{
        out.pop_back();
    }}
    if (!out.empty() && out.back() == '.') {{
        out.pop_back();
    }}
    return out;
}}

std::string escapeString(const std::string& input) {{
    std::string out;
    out.push_back('"');
    for (char c : input) {{
        if (c == '\\\\' || c == '"') {{
            out.push_back('\\\\');
            out.push_back(c);
        }} else if (c == '\\n') {{
            out += "\\\\n";
        }} else if (c == '\\r') {{
            out += "\\\\r";
        }} else if (c == '\\t') {{
            out += "\\\\t";
        }} else {{
            out.push_back(c);
        }}
    }}
    out.push_back('"');
    return out;
}}

std::string render(const Value& value);

std::string render(const Value& value) {{
    if (std::holds_alternative<std::nullptr_t>(value.data)) {{
        return "null";
    }}
    if (auto ptr = std::get_if<bool>(&value.data)) {{
        return *ptr ? "true" : "false";
    }}
    if (auto ptr = std::get_if<double>(&value.data)) {{
        return numberToString(*ptr);
    }}
    if (auto ptr = std::get_if<std::string>(&value.data)) {{
        return escapeString(*ptr);
    }}
    if (auto ptr = std::get_if<Value::Array>(&value.data)) {{
        std::string out = "[";
        for (size_t i = 0; i < ptr->size(); ++i) {{
            if (i > 0) {{
                out += ",";
            }}
            out += render((*ptr)[i]);
        }}
        out += "]";
        return out;
    }}
    const auto& obj = std::get<Value::Object>(value.data);
    std::string out = "{{";
    bool first = true;
    for (const auto& item : obj) {{
        if (!first) {{
            out += ",";
        }}
        first = false;
        out += escapeString(item.first);
        out += ":";
        out += render(item.second);
    }}
    out += "}}";
    return out;
}}

bool equalsValue(const Value& a, const Value& b) {{
    if (a.data.index() != b.data.index()) {{
        return false;
    }}
    if (std::holds_alternative<std::nullptr_t>(a.data)) {{
        return true;
    }}
    if (auto pa = std::get_if<bool>(&a.data)) {{
        return *pa == std::get<bool>(b.data);
    }}
    if (auto pa = std::get_if<double>(&a.data)) {{
        return std::fabs(*pa - std::get<double>(b.data)) < 1e-9;
    }}
    if (auto pa = std::get_if<std::string>(&a.data)) {{
        return *pa == std::get<std::string>(b.data);
    }}
    if (auto pa = std::get_if<Value::Array>(&a.data)) {{
        const auto& pb = std::get<Value::Array>(b.data);
        if (pa->size() != pb.size()) {{
            return false;
        }}
        for (size_t i = 0; i < pa->size(); ++i) {{
            if (!equalsValue((*pa)[i], pb[i])) {{
                return false;
            }}
        }}
        return true;
    }}
    const auto& oa = std::get<Value::Object>(a.data);
    const auto& ob = std::get<Value::Object>(b.data);
    if (oa.size() != ob.size()) {{
        return false;
    }}
    for (const auto& item : oa) {{
        auto it = ob.find(item.first);
        if (it == ob.end()) {{
            return false;
        }}
        if (!equalsValue(item.second, it->second)) {{
            return false;
        }}
    }}
    return true;
}}

Value sortPrimitiveArray(const Value::Array& arr) {{
    Value::Array out = arr;
    std::sort(out.begin(), out.end(), [](const Value& left, const Value& right) {{
        return render(left) < render(right);
    }});
    return Value(out);
}}

Value canonical(const Value& v, const std::string& mode) {{
    if (mode == "pair_unordered" || mode == "list_unordered" || mode == "strings_unordered") {{
        if (!std::holds_alternative<Value::Array>(v.data)) {{
            return v;
        }}
        return sortPrimitiveArray(std::get<Value::Array>(v.data));
    }}
    if (mode == "groups_unordered" || mode == "nested_unordered") {{
        if (!std::holds_alternative<Value::Array>(v.data)) {{
            return v;
        }}
        Value::Array outer;
        for (const Value& item : std::get<Value::Array>(v.data)) {{
            if (std::holds_alternative<Value::Array>(item.data)) {{
                outer.push_back(sortPrimitiveArray(std::get<Value::Array>(item.data)));
            }} else {{
                outer.push_back(item);
            }}
        }}
        std::sort(outer.begin(), outer.end(), [](const Value& left, const Value& right) {{
            return render(left) < render(right);
        }});
        return Value(outer);
    }}
    return v;
}}

bool equalFloatSequence(const Value& got, const Value& expected) {{
    if (!std::holds_alternative<Value::Array>(got.data) || !std::holds_alternative<Value::Array>(expected.data)) {{
        return false;
    }}
    const auto& ga = std::get<Value::Array>(got.data);
    const auto& ea = std::get<Value::Array>(expected.data);
    if (ga.size() != ea.size()) {{
        return false;
    }}
    for (size_t i = 0; i < ga.size(); ++i) {{
        if (ga[i].isNull() && ea[i].isNull()) {{
            continue;
        }}
        if (ga[i].isNull() || ea[i].isNull()) {{
            return false;
        }}
        if (std::fabs(toDoubleValue(ga[i]) - toDoubleValue(ea[i])) > 1e-9) {{
            return false;
        }}
    }}
    return true;
}}

bool equalValues(const Value& got, const Value& expected, const std::string& mode) {{
    if (mode == "pair_target_1idx") {{
        return false;
    }}
    if (mode == "float_sequence") {{
        return equalFloatSequence(got, expected);
    }}
    Value cg = canonical(got, mode);
    Value ce = canonical(expected, mode);
    if (equalsValue(cg, ce)) {{
        return true;
    }}
    if (cg.isNull() && std::holds_alternative<Value::Array>(ce.data) && std::get<Value::Array>(ce.data).empty()) {{
        return true;
    }}
    if (ce.isNull() && std::holds_alternative<Value::Array>(cg.data) && std::get<Value::Array>(cg.data).empty()) {{
        return true;
    }}
    return false;
}}

double toDoubleValue(const Value& v) {{
    if (auto ptr = std::get_if<double>(&v.data)) {{
        return *ptr;
    }}
    if (auto ptr = std::get_if<bool>(&v.data)) {{
        return *ptr ? 1.0 : 0.0;
    }}
    return 0.0;
}}

int toInt(const Value& v) {{
    return static_cast<int>(std::llround(toDoubleValue(v)));
}}

std::string toStringValue(const Value& v) {{
    if (auto ptr = std::get_if<std::string>(&v.data)) {{
        return *ptr;
    }}
    if (auto ptr = std::get_if<double>(&v.data)) {{
        return numberToString(*ptr);
    }}
    if (auto ptr = std::get_if<bool>(&v.data)) {{
        return *ptr ? "true" : "false";
    }}
    return "";
}}

bool toBoolValue(const Value& v) {{
    if (auto ptr = std::get_if<bool>(&v.data)) {{
        return *ptr;
    }}
    return false;
}}

std::vector<int> toIntVector(const Value& v) {{
    std::vector<int> out;
    for (const Value& item : asArray(v)) {{
        out.push_back(toInt(item));
    }}
    return out;
}}

std::vector<std::vector<int>> toIntMatrix(const Value& v) {{
    std::vector<std::vector<int>> out;
    for (const Value& row : asArray(v)) {{
        out.push_back(toIntVector(row));
    }}
    return out;
}}

std::vector<std::string> toStringVector(const Value& v) {{
    std::vector<std::string> out;
    for (const Value& item : asArray(v)) {{
        out.push_back(toStringValue(item));
    }}
    return out;
}}

std::vector<std::vector<std::string>> toStringMatrix(const Value& v) {{
    std::vector<std::vector<std::string>> out;
    for (const Value& row : asArray(v)) {{
        out.push_back(toStringVector(row));
    }}
    return out;
}}

std::vector<char> toCharVector(const Value& v) {{
    std::vector<char> out;
    for (const Value& item : asArray(v)) {{
        std::string s = toStringValue(item);
        out.push_back(s.empty() ? '\\0' : s[0]);
    }}
    return out;
}}

std::vector<std::vector<char>> toCharMatrix(const Value& v) {{
    std::vector<std::vector<char>> out;
    for (const Value& row : asArray(v)) {{
        out.push_back(toCharVector(row));
    }}
    return out;
}}

Value getArg(const Value::Object& testObj, const std::vector<std::string>& argKeys, size_t idx) {{
    auto itInput = testObj.find("input");
    if (itInput != testObj.end()) {{
        const Value& input = itInput->second;
        if (std::holds_alternative<Value::Object>(input.data)) {{
            Value::Object inputObj = std::get<Value::Object>(input.data);
            if (idx < argKeys.size()) {{
                return getField(inputObj, argKeys[idx]);
            }}
            return Value();
        }}
        if (std::holds_alternative<Value::Array>(input.data)) {{
            const auto& arr = std::get<Value::Array>(input.data);
            if (argKeys.size() > 1) {{
                if (idx == 0) {{
                    return input;
                }}
                if (idx < argKeys.size()) {{
                    auto it = testObj.find(argKeys[idx]);
                    if (it != testObj.end()) {{
                        return it->second;
                    }}
                }}
                if (idx < arr.size()) {{
                    return arr[idx];
                }}
                return Value();
            }}
            return input;
        }}
        return input;
    }}
    if (idx < argKeys.size()) {{
        auto it = testObj.find(argKeys[idx]);
        if (it != testObj.end()) {{
            return it->second;
        }}
    }}
    return Value();
}}

{list_helpers}
{tree_helpers}
{graph_helpers}

int main(int argc, char** argv) {{
    std::string testsPath = "../tests.yaml";
    std::string challengePath = "../challenge.yaml";
    if (argc > 2) {{
        testsPath = argv[2];
    }}
    if (argc > 3) {{
        challengePath = argv[3];
    }}
    std::ifstream testsFile(testsPath);
    std::ifstream challengeFile(challengePath);
    (void)testsFile;
    (void)challengeFile;

    std::vector<Value> tests = {{
        {tests_literal}
    }};
    std::vector<std::string> argKeys = {{{arg_keys_literal}}};
    bool failed = false;

    for (size_t i = 0; i < tests.size(); ++i) {{
        Value::Object testObj = asObject(tests[i]);
        Value gotValue;
{call_block}
        Value expected = getField(testObj, "expected");
        bool hasExpectedLength = hasField(testObj, "expected_length");
        bool passed = false;
        if (hasExpectedLength) {{
            Value expectedLength = getField(testObj, "expected_length");
            if (std::holds_alternative<Value::Array>(gotValue.data)) {{
                passed = static_cast<int>(std::get<Value::Array>(gotValue.data).size()) == toInt(expectedLength);
            }} else {{
                passed = false;
            }}
        }} else if ({json.dumps(mode)} == std::string("pair_target_1idx")) {{
            Value::Object input = asObject(getField(testObj, "input"));
            Value::Array numbers = asArray(getField(input, "numbers"));
            int target = toInt(getField(input, "target"));
            if (std::holds_alternative<Value::Array>(gotValue.data)) {{
                const Value::Array& pair = std::get<Value::Array>(gotValue.data);
                if (pair.size() == 2) {{
                    int i1 = toInt(pair[0]) - 1;
                    int i2 = toInt(pair[1]) - 1;
                    if (i1 >= 0 && i2 >= 0 && i1 < static_cast<int>(numbers.size()) && i2 < static_cast<int>(numbers.size()) && i1 != i2) {{
                        passed = toInt(numbers[i1]) + toInt(numbers[i2]) == target;
                    }} else {{
                        passed = false;
                    }}
                }} else {{
                    passed = false;
                }}
            }} else {{
                passed = false;
            }}
        }} else {{
            passed = equalValues(gotValue, expected, {json.dumps(mode)});
        }}

        int idx = static_cast<int>(i + 1);
        if (passed) {{
            std::cout << "PASS " << idx << "\\n";
        }} else {{
            failed = true;
            if (hasExpectedLength) {{
                std::cout << "FAIL " << idx << " got=" << render(gotValue) << " expected=" << render(getField(testObj, "expected_length")) << "\\n";
            }} else {{
                std::cout << "FAIL " << idx << " got=" << render(gotValue) << " expected=" << render(expected) << "\\n";
            }}
        }}
    }}

    if (failed) {{
        return 1;
    }}
    return 0;
}}
"""


def rust_builder(challenge, tests, funcs, impl_methods, arg_keys):
    return "fn main() { eprintln!(\"builder not implemented\"); std::process::exit(1); }\n"

def ensure_builder_dir(challenge):
    builder_dir = CHALLENGES_DIR / challenge / "builder"
    builder_dir.mkdir(parents=True, exist_ok=True)
    return builder_dir


def main():
    challenges = discover_challenges()
    for challenge in challenges:
        challenge_dir = CHALLENGES_DIR / challenge
        tests = normalize_tests(load_yaml(challenge_dir / "tests.yaml"))
        py_funcs, py_classes = parse_py_signature(challenge_dir / "setup" / "python.py")
        js_funcs, js_classes = parse_js_signature(challenge_dir / "setup" / "javascript.js")
        go_funcs, go_methods = parse_go_signature(challenge_dir / "setup" / "go.go")
        cpp_solution_methods, cpp_class_methods = parse_cpp_signature(challenge_dir / "setup" / "cpp.cpp")
        arg_keys = arg_keys_from_tests(tests)
        if not arg_keys and py_funcs:
            first_name = next(iter(py_funcs.keys()))
            arg_keys = list(py_funcs.get(first_name, []))
        builder_dir = ensure_builder_dir(challenge)
        (builder_dir / "builder.py").write_text(py_builder(challenge, tests, py_funcs, py_classes, arg_keys), encoding="utf-8")
        (builder_dir / "builder.js").write_text(js_builder(challenge, tests, js_funcs, js_classes, arg_keys), encoding="utf-8")
        (builder_dir / "builder.go").write_text(go_builder(challenge, tests, go_funcs, go_methods, arg_keys), encoding="utf-8")
        (builder_dir / "builder.cpp").write_text(cpp_builder(challenge, tests, cpp_solution_methods, cpp_class_methods, arg_keys), encoding="utf-8")
        (builder_dir / "builder.rs").write_text(rust_builder(challenge, tests, {}, {}, arg_keys), encoding="utf-8")


if __name__ == "__main__":
    main()
