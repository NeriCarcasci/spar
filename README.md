<p align="center">
  <img src="assets/banner.png" alt="spar" width="800"/>
</p>

<p align="center">
  <code>go install github.com/spar-cli/spar@latest</code>
</p>

<p align="center">
  <a href="#install"><img src="https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go 1.22+"/></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-333333?style=flat-square" alt="MIT License"/></a>
  <a href="#contributing"><img src="https://img.shields.io/badge/PRs-welcome-FF3B30?style=flat-square" alt="PRs Welcome"/></a>
</p>

---

spar is a free, open-source, terminal-based coding challenge platform. It has a built-in editor that blocks copy-paste, an AI interviewer that pressure-tests your solutions mid-session, and zero backend — everything runs locally.

The interview doesn't have a paste button. Neither does spar.

## Why

LeetCode is bloated, paywalled, and runs in a browser where pasting from ChatGPT takes two seconds. That's not practice — that's performance. spar puts you in your terminal with nothing but a problem and a blinking cursor. Every character comes from your fingers. Every answer comes from your head.

## How it works

```
$ spar
```

That's it. Pick a challenge. Write your solution in the built-in editor. Run the tests. If you've configured an AI provider, it interviews you about your approach — time complexity, edge cases, trade-offs — just like a real technical screen.

No browser. No account. No paywall.

## Features

**Built-in editor** — Syntax highlighting, line numbers, undo/redo, auto-indent. Paste detection rejects clipboard input and logs the attempt. Type it out.

**AI interviewer** — Three modes. *Interview mode* asks probing questions mid-solve ("What's the time complexity?", "What happens with empty input?"). *Practice mode* gives Socratic hints without ever revealing the answer. *Post-mortem mode* walks you through optimal solutions after you submit.

**165 challenges across 8 collections** — The Foundation (75 core algorithm patterns), System Design Lite, Concurrency & Parallelism, Data Structures from Scratch, Bit Manipulation & Math, Recursion & Backtracking, Real-World Patterns, and Language Idiomatic problems.

**5 languages** — Python, Go, JavaScript, C++, Rust. Every challenge has idiomatic setup code and reference solutions in all five.

**Mock interviews** — Curated sets of 3 problems designed to simulate a 45-minute technical screen. The AI treats the entire set as one continuous session.

**Your AI, your cost** — spar connects to your existing Claude or OpenAI account via OAuth. No API keys. No subscription to us. No middleman.

<!-- 
Uncomment when screenshots are ready:

## Screenshots

<p align="center">
  <img src="assets/screenshot-dashboard.png" alt="Dashboard" width="700"/>
</p>

<p align="center">
  <img src="assets/screenshot-session.png" alt="Coding session" width="700"/>
</p>
-->

## Install

**Go install** (recommended):

```bash
go install github.com/spar-cli/spar@latest
```

**From source:**

```bash
git clone https://github.com/spar-cli/spar.git
cd spar
go build -o spar ./cmd/spar
```

spar expects the challenge repo to be cloned locally. On first run it will ask you to configure the path.

### Language support

spar runs your solution using whatever toolchains you have installed. You need at least one:

| Language | Requirement |
|----------|------------|
| Python | `python3` on PATH |
| Go | `go` on PATH |
| JavaScript | `node` on PATH |
| C++ | `g++` on PATH |
| Rust | `rustc` on PATH |

No language installed? spar tells you. It doesn't guess.

## Challenges

Challenges live in `challenges/`, organized by collection and category. Each challenge folder contains:

```
challenges/arrays/two-sum/
  challenge.yaml          Problem description, constraints, hints
  tests.yaml              Visible + hidden test cases
  setup/                  Blank starting code per language
  solutions/              Reference solutions per language
```

The full challenge index is in `challenges/index.yaml` — a single file the app reads on startup for fast browsing. It's auto-generated, never hand-edited.

### Collections

| Collection | Count | What it covers |
|-----------|-------|----------------|
| The Foundation | 75 | Core algorithm patterns — arrays, trees, graphs, DP, the works |
| System Design Lite | 15 | LRU caches, rate limiters, pub/sub, circuit breakers |
| Concurrency | 10 | Producer-consumer, dining philosophers, deadlock detection |
| Data Structures | 15 | Build it yourself — hash maps, heaps, AVL trees, no stdlib |
| Bit Manipulation | 15 | The category everyone skips until the interview |
| Recursion Deep Dive | 10 | Constraint satisfaction, pruning, advanced backtracking |
| Real-World Patterns | 15 | JSON parsers, cron expressions, CLI arg parsing, diff algorithms |
| Language Idiomatic | 10 | Same problem, fundamentally different solution per language |

Plus 10 **mock interview sets** — curated triplets that simulate a 45-minute screen.

## Configuration

spar stores config in `~/.config/spar/config.yaml` and user data in `~/.local/share/spar/`. Created automatically on first run.

```yaml
# ~/.config/spar/config.yaml
repo_path: ~/code/spar
default_language: go
ai_provider: claude          # claude | openai | none
editor_tab_width: 4
```

## CLI

```
spar                    Launch the TUI
spar generate-index     Regenerate index.yaml from challenge tree
spar validate           Validate all challenge folders
spar validate <path>    Validate a specific challenge
spar version            Print version
```

## Contributing

Contributions are welcome — especially new challenges.

### Adding a challenge

1. Fork the repo
2. Create `challenges/{category}/{your-challenge}/` with all required files
3. Run `spar validate challenges/{category}/{your-challenge}` to check structure
4. Run `spar generate-index` to update the manifest
5. Open a PR

CI validates everything: correct structure, all declared languages have setup and solution files, all solutions pass all tests. If CI is green, you're good.

### Challenge quality bar

Every challenge needs original problem descriptions (not copied from other platforms), at least 2 visible and 3 hidden test cases covering edge cases, idiomatic solutions in all supported languages (not transliterated Python), and no comments in code files — the code speaks for itself.

## Philosophy

spar is built on a few beliefs:

Coding skill comes from writing code, not reading solutions. The terminal is where developers actually work — meet them there. Copy-paste is the enemy of learning. AI should be an interviewer, not an answer key. Open source means free forever, not free-until-we-raise-a-Series-A.

## License

[MIT](LICENSE)

---

<p align="center">
  <sub>code under pressure.</sub>
</p>
