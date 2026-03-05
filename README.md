# clog

A minimal CLI tool for developers who want to see what they actually shipped, across multiple repos, without opening GitHub.

Run `clog` at the end of the day to see a clean list of your commits. Run `clog sum` to get an AI-generated summary ready to paste into your standup, PR description, or team update.

---

## Features

- Shows today's or this week's commits across all your tracked repos
- Groups output cleanly by repo
- Filters to your commits only
- AI-powered summaries via Claude (optional)
- Zero external dependencies - pure Go, single binary

---

## Install

```bash
git clone https://github.com/atkntepe/clog
cd clog
go build -o clog .
cp clog /usr/local/bin/clog
```

Requires Go 1.25+.

---

## Setup

**Add repos to track:**
```bash
clog repo --add frontend ~/projects/my-frontend
clog repo --add api ~/projects/my-api
clog repo --add mobile ~/projects/my-mobile
```

**Set your name** (used to filter to your commits only):
```bash
clog config --author "Your Name"
```

**For AI summaries**, set your Anthropic API key:
```bash
clog config --api-key "your-key"
```

**Optionally set a different model** (defaults to `claude-haiku-4-5-20251001`):
```bash
clog config --model "claude-haiku-4-5-20251001"
```

---

## Usage

```bash
clog               # today's commits across all repos
clog week          # this week's commits
clog sum           # today's commits + AI standup summary
clog sum --week    # this week's commits + AI summary
```

**Managing repos:**
```bash
clog repo --list              # show all tracked repos
clog repo --add name /path    # add a repo
clog repo --remove name       # remove a repo
```

**Configuration:**
```bash
clog config --author "Your Name"       # set git author name
clog config --api-key "your-key"       # set Anthropic API key
clog config --model "model-name"       # set AI model
```

---

## Example Output

```
● frontend
  a3f1c2e  fix: resolve form validation on Safari
  b7d09f1  feat: add dark mode toggle to settings page

● api
  cc4812a  fix: correct rate limit headers on auth endpoints
  d901ee3  chore: update dependencies

  4 commits across 2 repos
```

With `clog sum`:

```
● frontend
  a3f1c2e  fix: resolve form validation on Safari
  b7d09f1  feat: add dark mode toggle to settings page

● api
  cc4812a  fix: correct rate limit headers on auth endpoints
  d901ee3  chore: update dependencies

  4 commits across 2 repos

─────────────────────────────────
  AI Summary

  Today I fixed a Safari form validation bug and shipped dark mode support on the
  settings page. On the API side, I corrected rate limit response headers and kept
  dependencies up to date.
─────────────────────────────────
```

---

## Configuration

Config is stored at `~/.config/clog/config.json`. It contains repo paths, your author name, and the model setting.

```json
{
  "author": "Your Name",
  "model": "claude-haiku-4-5-20251001",
  "repos": [
    { "name": "frontend", "path": "/Users/you/projects/frontend" },
    { "name": "api", "path": "/Users/you/projects/api" }
  ]
}
```

The API key is stored separately in `~/.config/clog/.env` with restricted permissions (owner-only). Environment variables (`ANTHROPIC_API_KEY`, `ANTHROPIC_MODEL`) override stored values if set.

---

## Why clog?

Most days you open your standup message and stare at a blank field trying to remember what you did. `clog` solves that. It's a one-command answer to *"what did I ship today?"*

---

## License

MIT
