---
name: go-code-review
description: >-
  Review Go code for idiomatic quality, maintainability, and testability in
  cards-cli. Use when reviewing Go changes, pull requests, refactors, or when
  the user asks for a Go code review, best practices check, or feedback on
  error handling, duplication, function size, or pure helper extraction.
---

# Go Code Review (cards-cli)

Review Go code against project conventions and idiomatic Go practices. Focus on maintainability and testability, not style nitpicks.

## Review workflow

1. Read the diff or changed files end-to-end before commenting.
2. Run `go test ./...` and `go vet ./...` when reviewing local changes.
3. Apply the checklist below; cite specific files and lines.
4. Separate **must fix** issues from **suggestions**.

## Checklist

### Single responsibility and function size

- [ ] Each function does one thing. Flag handlers or repository methods that mix validation, I/O, formatting, and orchestration.
- [ ] Prefer short functions (~30–50 lines as a soft ceiling). Extract steps into named helpers when logic branches or nests deeply.
- [ ] Cobra `RunE` handlers should orchestrate only: parse flags → call pure/domain logic → call I/O → format output.

### Pure helpers (testability)

- [ ] **Validation, parsing, formatting, and mapping** should be pure functions: same inputs → same outputs, no globals, no I/O, no `time.Now()` hidden inside (inject or pass timestamps).
- [ ] Extract pure logic into `internal/models`, `internal/output`, or package-level unexported helpers with `_test.go` coverage.
- [ ] Side effects (DB, filesystem, stdout) stay at the edges; core logic returns `(T, error)` or `error`.

### Avoid duplication

- [ ] Repeated validation, SQL fragments, flag parsing, or error wrapping → extract a shared helper.
- [ ] Similar commands sharing patterns → reuse `runWithRepo`, model helpers, and formatters instead of copy-pasting.

### Error handling

- [ ] Errors are checked; never `_ = err` or silent ignores.
- [ ] Wrap errors with context: `fmt.Errorf("insert card: %w", err)`. Use `%w` when callers need `errors.Is` / `errors.As`.
- [ ] Sentinel errors for domain cases — compare with `errors.Is`, do not string-match.
- [ ] Return errors; do not `log.Fatal` / `os.Exit` outside `main`.
- [ ] User-facing messages belong at the CLI boundary; internal packages return descriptive but neutral errors.

### Interfaces and dependencies

- [ ] Depend on small interfaces (`io.Writer`, `Formatter`) at call sites that need test doubles.
- [ ] Accept `context.Context` as the first parameter on I/O-bound functions (DB, HTTP).
- [ ] Constructor injection (`NewRepository(db *DB)`) over package-level mutable state.

### Naming and packages

- [ ] Package names are short, lowercase, single-word (`db`, `models`, `output`).
- [ ] Avoid stutter: `models.Card` not `models.CardModel`.
- [ ] Receivers: one or two letters from the type name (`c *Card`, `r *Repository`).

### Tests

- [ ] Pure helpers have table-driven unit tests with edge cases.
- [ ] I/O code tested with fakes or `bytes.Buffer` / in-memory SQLite where practical.
- [ ] Test names: `TestValidateFront` / subtests `t.Run("empty", ...)`.
- [ ] No logic in tests that duplicates production code without asserting behavior.

## Severity labels

| Label | When |
|-------|------|
| **Critical** | Bugs, data loss, unchecked errors, broken API contracts |
| **Should fix** | Duplication, untestable logic, poor error wrapping, SRP violations |
| **Suggestion** | Naming, minor simplification, optional refactor |
| **Nice to have** | Comments, micro-optimizations |

## Commands

```bash
go test ./...
go vet ./...
go test -race ./...   # when concurrency is involved
```
