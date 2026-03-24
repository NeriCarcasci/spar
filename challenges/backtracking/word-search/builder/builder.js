const fs = require("fs");
const path = require("path");

const TESTS = JSON.parse("[{\"input\":{\"board\":[[\"A\",\"B\",\"C\",\"E\"],[\"S\",\"F\",\"C\",\"S\"],[\"A\",\"D\",\"E\",\"E\"]],\"word\":\"ABCCED\"},\"expected\":true,\"visible\":true},{\"input\":{\"board\":[[\"A\",\"B\",\"C\",\"E\"],[\"S\",\"F\",\"C\",\"S\"],[\"A\",\"D\",\"E\",\"E\"]],\"word\":\"SEE\"},\"expected\":true,\"visible\":true},{\"input\":{\"board\":[[\"A\",\"B\",\"C\",\"E\"],[\"S\",\"F\",\"C\",\"S\"],[\"A\",\"D\",\"E\",\"E\"]],\"word\":\"ABCB\"},\"expected\":false,\"visible\":true},{\"input\":{\"board\":[[\"A\"]],\"word\":\"A\"},\"expected\":true,\"visible\":false},{\"input\":{\"board\":[[\"A\",\"B\"],[\"C\",\"D\"]],\"word\":\"ABCDE\"},\"expected\":false,\"visible\":false},{\"input\":{\"board\":[[\"A\",\"A\",\"A\"],[\"A\",\"A\",\"A\"],[\"A\",\"A\",\"A\"]],\"word\":\"AAAAAAAAA\"},\"expected\":true,\"visible\":false},{\"input\":{\"board\":[[\"A\",\"B\"],[\"C\",\"D\"]],\"word\":\"ACDB\"},\"expected\":true,\"visible\":false}]");
const COMPARE_MODE = "exact";

function isObject(v) {
    return v !== null && typeof v === "object" && !Array.isArray(v);
}

function resolveFunction(solution, names) {
    for (const name of names) {
        if (typeof solution[name] === "function") {
            return solution[name];
        }
    }
    throw new Error("missing function");
}

function resolveClass(solution, name, fallback) {
    if (typeof solution[name] === "function") {
        return solution[name];
    }
    return fallback;
}

class ListNodeFallback {
    constructor(val = 0, next = null) {
        this.val = val;
        this.next = next;
    }
}

class TreeNodeFallback {
    constructor(val = 0, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

class GraphNodeFallback {
    constructor(val = 0, neighbors = []) {
        this.val = val;
        this.neighbors = neighbors;
    }
}

function arrayToLinkedList(values, NodeClass) {
    let head = null;
    let tail = null;
    for (const value of values) {
        const node = new NodeClass(value);
        if (head === null) {
            head = node;
            tail = node;
        } else {
            tail.next = node;
            tail = node;
        }
    }
    return head;
}

function linkedListToArray(head) {
    const out = [];
    let node = head;
    let guard = 0;
    while (node !== null && node !== undefined && guard < 100000) {
        out.push(node.val);
        node = node.next;
        guard += 1;
    }
    return out;
}

function arrayToTree(values, TreeNodeClass) {
    if (!Array.isArray(values) || values.length === 0 || values[0] === null) {
        return null;
    }
    const nodes = values.map((v) => (v === null ? null : new TreeNodeClass(v)));
    let idx = 1;
    for (const node of nodes) {
        if (node === null) {
            continue;
        }
        if (idx < nodes.length) {
            node.left = nodes[idx];
            idx += 1;
        }
        if (idx < nodes.length) {
            node.right = nodes[idx];
            idx += 1;
        }
    }
    return nodes[0];
}

function treeToArray(root) {
    if (root === null || root === undefined) {
        return [];
    }
    const out = [];
    const queue = [root];
    let i = 0;
    while (i < queue.length) {
        const node = queue[i];
        i += 1;
        if (node === null) {
            out.push(null);
            continue;
        }
        out.push(node.val);
        queue.push(node.left === undefined ? null : node.left);
        queue.push(node.right === undefined ? null : node.right);
    }
    while (out.length > 0 && out[out.length - 1] === null) {
        out.pop();
    }
    return out;
}

function adjListToGraph(adjList, NodeClass) {
    if (!Array.isArray(adjList) || adjList.length === 0) {
        return null;
    }
    const nodes = [];
    for (let i = 0; i < adjList.length; i += 1) {
        nodes.push(new NodeClass(i + 1, []));
    }
    for (let i = 0; i < adjList.length; i += 1) {
        nodes[i].neighbors = (adjList[i] || []).map((n) => nodes[n - 1]);
    }
    return nodes[0];
}

function graphToAdjList(root) {
    if (root === null || root === undefined) {
        return [];
    }
    const seen = new Map();
    const queue = [root];
    while (queue.length > 0) {
        const node = queue.shift();
        if (seen.has(node.val)) {
            continue;
        }
        seen.set(node.val, node);
        for (const nei of node.neighbors || []) {
            queue.push(nei);
        }
    }
    let maxId = 0;
    for (const id of seen.keys()) {
        if (id > maxId) {
            maxId = id;
        }
    }
    const out = Array.from({ length: maxId }, () => []);
    for (const [id, node] of seen.entries()) {
        const neighbors = (node.neighbors || []).map((n) => n.val).sort((a, b) => a - b);
        out[id - 1] = neighbors;
    }
    return out;
}

function canonical(value, mode) {
    if (mode === "pair_unordered" || mode === "list_unordered" || mode === "strings_unordered") {
        if (Array.isArray(value)) {
            return [...value].sort();
        }
        return value;
    }
    if (mode === "groups_unordered") {
        if (!Array.isArray(value)) {
            return value;
        }
        return value.map((group) => (Array.isArray(group) ? [...group].sort() : group)).sort((a, b) => JSON.stringify(a).localeCompare(JSON.stringify(b)));
    }
    if (mode === "nested_unordered") {
        if (!Array.isArray(value)) {
            return value;
        }
        return value.map((item) => (Array.isArray(item) ? [...item].sort((a, b) => (a > b ? 1 : a < b ? -1 : 0)) : item)).sort((a, b) => JSON.stringify(a).localeCompare(JSON.stringify(b)));
    }
    return value;
}

function equalValues(got, expected, mode) {
    if (mode === "pair_target_1idx") {
        return false;
    }
    if (mode === "float_sequence") {
        if (!Array.isArray(got) || !Array.isArray(expected) || got.length !== expected.length) {
            return false;
        }
        for (let i = 0; i < got.length; i += 1) {
            const a = got[i];
            const b = expected[i];
            if (a === null && b === null) {
                continue;
            }
            if (a === null || b === null) {
                return false;
            }
            if (Math.abs(Number(a) - Number(b)) > 1e-9) {
                return false;
            }
        }
        return true;
    }
    return JSON.stringify(canonical(got, mode)) === JSON.stringify(canonical(expected, mode));
}

function render(value) {
    return JSON.stringify(value);
}

function loadSolution(solutionPath) {
    const resolved = path.resolve(solutionPath);
    delete require.cache[resolved];
    return require(resolved);
}

function run(solutionPath, testsPath, challengePath) {
    fs.readFileSync(testsPath, "utf-8");
    fs.readFileSync(challengePath, "utf-8");
    const solution = loadSolution(solutionPath);
    let failures = 0;
    for (let i = 0; i < TESTS.length; i += 1) {
        const test = TESTS[i];
        let gotValue;

    const fn = resolveFunction(solution, ["exist"]);
    const caseInput = test.input;
    let args = [];
    if (isObject(caseInput)) {
        args = ["board", "word"].map((key) => caseInput[key]);
    } else if (Array.isArray(caseInput) && ["board", "word"].length > 1) {
        args.push(caseInput);
        for (let i = 1; i < ["board", "word"].length; i += 1) {
            const key = ["board", "word"][i];
            if (Object.prototype.hasOwnProperty.call(test, key)) {
                args.push(test[key]);
            } else {
                args.push(i < caseInput.length ? caseInput[i] : null);
            }
        }
    } else if (caseInput !== undefined) {
        args = [caseInput];
    } else {
        args = ["board", "word"].map((key) => test[key]);
    }
    gotValue = fn(...args);

        const expected = test.expected;
        const expectedLength = test.expected_length;
        let passed = false;
        if (COMPARE_MODE === "pair_target_1idx") {
            const inp = test.input || {};
            const numbers = Array.isArray(inp.numbers) ? inp.numbers : [];
            const target = inp.target;
            if (Array.isArray(gotValue) && gotValue.length === 2 && target !== undefined && target !== null) {
                const i = Number(gotValue[0]) - 1;
                const j = Number(gotValue[1]) - 1;
                passed = i >= 0 && j >= 0 && i < numbers.length && j < numbers.length && i !== j && Number(numbers[i]) + Number(numbers[j]) === Number(target);
            } else {
                passed = false;
            }
        } else if (expectedLength !== undefined) {
            passed = Array.isArray(gotValue) && gotValue.length === expectedLength;
        } else {
            passed = equalValues(gotValue, expected, COMPARE_MODE);
        }
        const idx = i + 1;
        if (passed) {
            console.log(`PASS ${idx}`);
        } else {
            failures += 1;
            if (expectedLength !== undefined) {
                console.log(`FAIL ${idx} got=${render(gotValue)} expected=${render(expectedLength)}`);
            } else {
                console.log(`FAIL ${idx} got=${render(gotValue)} expected=${render(expected)}`);
            }
        }
    }
    return failures === 0;
}

function main() {
    const builderDir = __dirname;
    const solutionPath = process.argv[2] ? path.resolve(process.argv[2]) : path.resolve(builderDir, "..", "setup", "javascript.js");
    const testsPath = process.argv[3] ? path.resolve(process.argv[3]) : path.resolve(builderDir, "..", "tests.yaml");
    const challengePath = process.argv[4] ? path.resolve(process.argv[4]) : path.resolve(builderDir, "..", "challenge.yaml");
    const ok = run(solutionPath, testsPath, challengePath);
    process.exit(ok ? 0 : 1);
}

main();
