---
name: github-issues
description: >-
  Create GitHub issues for cards-cli with consistent structure, branch naming,
  acceptance criteria, and execution steps. Use when the user asks to create,
  draft, or open an issue, plan work as a GitHub issue, or file a bug/feature
  request with implementation and release steps.
---

# GitHub Issues (cards-cli)

Create issues via `gh` that match this repo's conventions. Read closed issues (`gh issue list --state all --limit 10`) when unsure of tone or structure.

## Workflow

1. Clarify goal, scope, and whether a **release** step is needed (user-facing changes → include release step).
2. Pick label: `enhancement` (features), `bug` (fixes), `documentation` (docs-only).
3. Draft the issue body from the template below.
4. Create the issue:

```bash
gh issue create --title "<title>" --label "<label>" --body "$(cat <<'EOF'
<issue body>
EOF
)"
```

5. Note the issue number from the URL. Update the body if it contains placeholders:
   - Branch name: `feat/<N>-short-description` (use actual number)
   - PR close line: `Closes #<N>`
   - Execution step 2: reference the final branch name

```bash
gh issue edit <N> --body "$(gh issue view <N> --json body -q .body | sed 's/Closes #<issue-number>/Closes #<N>/')"
```

## Branch naming

Pattern: `(feat|fix|doc)/<issue-number>-short-description`

| Prefix | Use for |
| --- | --- |
| `feat` | New features, enhancements, behavior changes |
| `fix` | Bug fixes |
| `doc` | Documentation-only work |

Use kebab-case after the number. Example: `feat/24-simplify-study-grades`.

## Issue template

```markdown
## Description

<Why this work matters — 1–3 sentences.>

## Branch

Create a feature branch from `main` using the naming pattern `(feat|fix|doc)/<issue-number>-description`:

```bash
git checkout main
git pull origin main
git checkout -b <prefix>/<issue-number>-short-description
```

## Acceptance criteria

### <area 1>
- <testable criterion>
- <testable criterion>

### <area 2>
- <testable criterion>

### Tests
- All tests pass: `go test ./...` and `go vet ./...`

### Documentation
- `CHANGELOG.md` — update `[Unreleased]` as needed
- `docs/COMMANDS.md`, `README.md`, `.cursor/skills/cards-cli/reference.md` — when commands or behavior change

## Files to change (expected)

| Area | Files |
| --- | --- |
| <area> | `<paths>` |

Omit this section for small or obvious changes.

## Execution steps

1. Assign issue to `joaovictornsv`
2. Create branch `<prefix>/<issue-number>-short-description` from latest `main` (see **Branch** section above)
3. Implement per acceptance criteria above
4. Run `go test ./...` and `go vet ./...`
5. Review the diff with `.agents/go-code-review/SKILL.md` — separate **must fix** from **suggestions**
6. Apply all review findings, including suggestions
7. Update `CHANGELOG.md`, `docs/COMMANDS.md`, `.cursor/skills/cards-cli/reference.md`, and `README.md` as needed
8. Commit, push, and open a PR (`Closes #<issue-number>`) with `--assignee joaovictornsv`
9. Merge PR with merge commit per `.agents/git-workflow.md`
10. **Release** (include when user-facing; per `.agents/github-releases/SKILL.md`):
    - Ensure CI is green on `main`
    - Move `[Unreleased]` changelog entry into a new version section (`## [x.y.z] - YYYY-MM-DD`)
    - Commit changelog: `chore: release vX.Y.Z`
    - Tag: `git tag -a vX.Y.Z -m vX.Y.Z`
    - Push `main` and tag
    - Build `cards-linux-amd64`, generate `SHA256SUMS`, create GitHub release with concise bullet notes
    - Attach `cards-linux-amd64` and `SHA256SUMS` to the release
```

## Section rules

| Section | Required | Notes |
| --- | --- | --- |
| Description | Yes | Problem/goal, not implementation details |
| Branch | Yes | Full git commands + naming pattern |
| Acceptance criteria | Yes | Testable bullets; group by area (UX, API, tests, docs) |
| Files to change | Optional | Helpful for multi-package changes |
| Execution steps | Yes | Always include assign → branch → implement → quality → PR → merge |
| Release step | When asked or user-facing | Omit for internal/CI-only work unless requested |

## Optional fields

Add at the top of the body when relevant:

```markdown
**Priority:** P<n>

**Dependencies:** #<issue> (<short note>)
```

## Title conventions

- Imperative, concise: `Simplify study session grades to again and easy`
- No issue number in title
- Bug issues may start with `Fix:` when it reads naturally

## Checklist before submitting

- [ ] Branch prefix matches work type (`feat` / `fix` / `doc`)
- [ ] Acceptance criteria are testable, not vague
- [ ] Execution steps reference `.agents/go-code-review/SKILL.md` and `.agents/git-workflow.md`
- [ ] Release step included when user asked or change is user-visible
- [ ] Placeholders replaced with actual issue number after creation
