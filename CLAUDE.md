# Damon — Nomad TUI

Damon is a terminal UI (Go, gocui-based) for observing and interacting with
HashiCorp Nomad resources — jobs, allocations, deployments, namespaces, task
events, and log streams. It talks to Nomad exclusively over its HTTP API via
the standard `NOMAD_*` environment variables (see README.md).

## Layout

- `cmd/` — entrypoint
- `view/` — TUI views/screens (gocui panels)
- `component/` — reusable UI components (tables, status, task events, errors)
- `watcher/` — polling/streaming logic against the Nomad API (jobs, allocations, logs, namespaces)
- `layout/` — overall screen layout/composition
- `nomad/` — Nomad API client wrapper
- `models/`, `primitives/`, `state/`, `styles/`, `refresher/` — supporting types, drawing primitives, shared state, styling, refresh scheduling

## Conventions

- Fakes under `*/*fakes/` (e.g. `view/viewfakes/`) are generated with
  [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) — don't hand-edit them,
  regenerate from the interface instead (`go generate ./...` if a `//go:generate` directive
  is present, otherwise re-run counterfeiter against the interface).
- Match existing test style in the package you're editing (most packages have a
  `_test.go` per file using table-driven tests).

## General guidelines

- When planning, verify against current version documentation rather than guessing.
- Check for updates to any relevant dependency before planning around it.

## Logging findings

When a non-obvious bug or behavior gets root-caused during a session —
something that cost real time to track down and would plausibly get
rediscovered the hard way again later (Nomad API quirks, gocui redraw/focus
gotchas, race conditions in the watcher/streaming code, etc.) — save it as a
memory entry rather than a file, so it surfaces automatically in future
sessions on this repo. Tell the user in the conversation when this fires.

Not every bug fix qualifies — routine typos or mistakes obvious from the diff
don't need one. Reserve it for things that took real investigation to
root-cause.

## Session log

Keep a running log of notable work done on this repo (what changed, why, and
any open threads) as memory entries rather than a `workbook.md` file, so it's
available across sessions without living in the repo.


## Versioning and Changelog

At the end of each fixed bug, or added feature, add a line to a changelog.md and bump the version in version/version.go.
