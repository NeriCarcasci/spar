import ast
import importlib.util
import json
import pathlib
import sys


def parse_value(raw: str):
    value = raw.strip()
    if value == "":
        return ""
    lowered = value.lower()
    if lowered == "null":
        return None
    if lowered == "true":
        return True
    if lowered == "false":
        return False
    if value.startswith("[") or value.startswith("{"):
        normalized = value.replace("null", "None").replace("true", "True").replace("false", "False")
        return ast.literal_eval(normalized)
    if (value.startswith('"') and value.endswith('"')) or (value.startswith("'") and value.endswith("'")):
        return value[1:-1]
    try:
        return int(value)
    except ValueError:
        pass
    try:
        return float(value)
    except ValueError:
        pass
    return value


def split_key_value(line: str):
    key, value = line.split(":", 1)
    return key.strip(), value.strip()


def parse_tests(path: pathlib.Path):
    lines = path.read_text(encoding="utf-8").splitlines()
    tests = []
    section = "visible"
    i = 0
    while i < len(lines):
        line = lines[i].rstrip("\n")
        stripped = line.strip()
        if not stripped:
            i += 1
            continue
        if stripped in ("visible:", "hidden:"):
            section = stripped[:-1]
            i += 1
            continue
        if line.startswith("  - "):
            case = {}
            first = line[4:].strip()
            i += 1
            if first:
                key, value = split_key_value(first)
                case[key] = parse_value(value)
            while i < len(lines):
                nxt = lines[i].rstrip("\n")
                nxt_stripped = nxt.strip()
                if nxt.startswith("  - ") or nxt_stripped in ("visible:", "hidden:"):
                    break
                if not nxt_stripped:
                    i += 1
                    continue
                indent = len(nxt) - len(nxt.lstrip(" "))
                if indent >= 4 and ":" in nxt:
                    key, value = split_key_value(nxt[indent:].strip())
                    case[key] = parse_value(value)
                i += 1
            case["visible"] = section == "visible"
            tests.append(case)
            continue
        i += 1
    return tests


def parse_challenge(path: pathlib.Path):
    meta = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        stripped = line.strip()
        if not stripped or stripped.startswith("#"):
            continue
        if line.startswith(" "):
            continue
        if ":" not in stripped:
            continue
        key, value = split_key_value(stripped)
        meta[key] = value.strip('"')
    return meta


def load_module(solution_path: pathlib.Path):
    spec = importlib.util.spec_from_file_location("user_solution", solution_path)
    if spec is None or spec.loader is None:
        raise RuntimeError("unable to load solution")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def resolve_function(module):
    for name in ("reverse_list", "reverseList", "ReverseList"):
        fn = getattr(module, name, None)
        if callable(fn):
            return fn
    raise RuntimeError("missing reverse_list function")


def resolve_list_node(module):
    cls = getattr(module, "ListNode", None)
    if cls is None:
        class ListNode:
            def __init__(self, val=0, next=None):
                self.val = val
                self.next = next
        return ListNode
    return cls


def array_to_linked_list(values, node_cls):
    head = None
    tail = None
    for value in values:
        node = node_cls(value)
        if head is None:
            head = node
            tail = node
            continue
        tail.next = node
        tail = node
    return head


def linked_list_to_array(head):
    values = []
    seen = 0
    node = head
    while node is not None and seen < 100000:
        values.append(node.val)
        node = node.next
        seen += 1
    return values


def render(value):
    if isinstance(value, (list, dict, tuple, int, float, bool)) or value is None:
        return json.dumps(value, separators=(",", ":"))
    return str(value)


def run(solution_path: pathlib.Path, tests_path: pathlib.Path, challenge_path: pathlib.Path):
    module = load_module(solution_path)
    fn = resolve_function(module)
    node_cls = resolve_list_node(module)
    tests = parse_tests(tests_path)
    meta = parse_challenge(challenge_path)
    input_type = meta.get("input_type", "array")
    output_type = meta.get("output_type", "array")

    failures = 0
    for idx, test in enumerate(tests, start=1):
        raw_input = test.get("input", [])
        if input_type == "linked-list":
            call_input = array_to_linked_list(raw_input, node_cls)
        else:
            call_input = raw_input

        got = fn(call_input)
        if output_type == "linked-list":
            got_value = linked_list_to_array(got)
        else:
            got_value = got

        expected = test.get("expected")
        if got_value == expected:
            print(f"PASS {idx}")
        else:
            failures += 1
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