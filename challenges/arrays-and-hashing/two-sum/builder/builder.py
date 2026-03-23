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
                if value:
                    case[key] = parse_value(value)
                else:
                    nested = {}
                    while i < len(lines):
                        sub = lines[i].rstrip("\n")
                        if not sub.strip():
                            i += 1
                            continue
                        indent = len(sub) - len(sub.lstrip(" "))
                        if indent <= 4:
                            break
                        sub_key, sub_value = split_key_value(sub.strip())
                        nested[sub_key] = parse_value(sub_value)
                        i += 1
                    case[key] = nested
            while i < len(lines):
                nxt = lines[i].rstrip("\n")
                nxt_stripped = nxt.strip()
                if nxt.startswith("  - ") or nxt_stripped in ("visible:", "hidden:"):
                    break
                if not nxt_stripped:
                    i += 1
                    continue
                indent = len(nxt) - len(nxt.lstrip(" "))
                if indent < 4:
                    i += 1
                    continue
                key, value = split_key_value(nxt[indent:].strip())
                if value:
                    case[key] = parse_value(value)
                    i += 1
                    continue
                nested = {}
                i += 1
                while i < len(lines):
                    sub = lines[i].rstrip("\n")
                    if not sub.strip():
                        i += 1
                        continue
                    sub_indent = len(sub) - len(sub.lstrip(" "))
                    if sub_indent <= indent:
                        break
                    sub_key, sub_value = split_key_value(sub[sub_indent:].strip())
                    nested[sub_key] = parse_value(sub_value)
                    i += 1
                case[key] = nested
            case["visible"] = section == "visible"
            tests.append(case)
            continue
        i += 1
    return tests


def load_module(solution_path: pathlib.Path):
    spec = importlib.util.spec_from_file_location("user_solution", solution_path)
    if spec is None or spec.loader is None:
        raise RuntimeError("unable to load solution")
    module = importlib.util.module_from_spec(spec)
    spec.loader.exec_module(module)
    return module


def resolve_function(module):
    for name in ("two_sum", "twoSum", "TwoSum"):
        fn = getattr(module, name, None)
        if callable(fn):
            return fn
    raise RuntimeError("missing two_sum function")


def render(value):
    if isinstance(value, (list, dict, tuple, int, float, bool)) or value is None:
        return json.dumps(value, separators=(",", ":"))
    return str(value)


def run(solution_path: pathlib.Path, tests_path: pathlib.Path):
    module = load_module(solution_path)
    fn = resolve_function(module)
    tests = parse_tests(tests_path)
    failures = 0
    for idx, test in enumerate(tests, start=1):
        raw_input = test.get("input")
        target = test.get("target")
        if isinstance(raw_input, dict):
            got = fn(**raw_input)
        elif target is not None:
            got = fn(raw_input, target)
        else:
            got = fn(raw_input)
        expected = test.get("expected")
        if got == expected:
            print(f"PASS {idx}")
        else:
            failures += 1
            print(f"FAIL {idx} got={render(got)} expected={render(expected)}")
    return failures == 0


def main():
    builder_dir = pathlib.Path(__file__).resolve().parent
    solution_path = pathlib.Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else (builder_dir.parent / "setup" / "python.py")
    tests_path = pathlib.Path(sys.argv[2]).resolve() if len(sys.argv) > 2 else (builder_dir.parent / "tests.yaml")
    ok = run(solution_path, tests_path)
    raise SystemExit(0 if ok else 1)


if __name__ == "__main__":
    main()