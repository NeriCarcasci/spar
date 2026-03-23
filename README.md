```
v
```

**terminal-native coding challenges with an AI interviewer**

---

## Install

```bash
go install github.com/spar/spar@latest
```

## What is spar?

spar is a free, open-source, terminal-based coding challenge platform. It has a built-in editor that blocks copy-paste, an AI interviewer that pressure-tests your solutions mid-session, and zero backend — everything runs locally.

The interview doesn't have a paste button. Neither does spar.

## How it works

Clone this repo. Run `spar`. Pick a challenge. Write your solution in the built-in editor. Run the tests. If you enabled an AI provider, it will interview you about your approach — asking about time complexity, edge cases, and trade-offs, just like a real technical interview.

No browser. No account. No paywall.

## Challenges

Challenges live in the `challenges/` directory, organized by category. Each challenge has:

- A problem description with examples and constraints
- Test cases (visible and hidden)
- Setup files for each supported language (Python, Go, JavaScript, C++, Rust)
- Reference solutions

Anyone can contribute challenges via pull request.

## Configuration

spar stores config in `~/.config/spar/config.yaml` and user data in `~/.local/share/spar/`. On first run it creates these directories automatically.

Set your challenge repo path, preferred language, and AI provider in the config file.

## CLI

```
spar                    Launch the TUI
spar generate-index     Regenerate index.yaml from challenge tree
spar validate           Validate all challenge folders
spar validate <path>    Validate a specific challenge folder
spar version            Print version
```

## Contributing

Contributions are welcome. Read the challenge format spec in the repo, write your challenge, run `spar validate` and `spar generate-index`, then open a PR.

## License

MIT
