# Git workflow rules

Follow these rules for all git operations in this repository.

## Prohibited operations

- **Never squash merge** — use a regular merge commit when merging pull requests.
- **Never rebase** — do not rebase branches onto `main` or other bases.
- **Never `git reset` to undo commits** — do not use reset to remove or rewrite committed history.
- **Never `git amend`** — do not amend existing commits; create a new commit instead.
- **Never force push** — do not use `git push --force` or `git push --force-with-lease`.

## Preferred workflow

1. Create a feature branch from `main`.
2. Commit changes on the branch (new commits only; no amend).
3. Open a pull request and wait for CI to pass.
4. Merge the PR with a **merge commit** (`gh pr merge --merge`).
5. Pull `main` locally after the merge.

## Release tags

- Tag release commits on `main` after merge.
- Push tags normally (`git push origin <tag>`); never force-push tags.
