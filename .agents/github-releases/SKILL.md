---
name: github-releases
description: >-
  Create and publish GitHub releases for cards-cli using semantic versioning
  and concise bullet-point release notes derived from CHANGELOG.md. Use when
  cutting a release, tagging a version, writing release notes, or updating
  CHANGELOG.md for a new version.
---

# GitHub Releases (cards-cli)

Simple release workflow for this project: semver tag, changelog notes, and a Linux amd64 binary attached to every GitHub release.

## Versioning

- Follow [Semantic Versioning](https://semver.org/): `MAJOR.MINOR.PATCH` (e.g. `v0.1.0`).
- Tag format: `v` prefix required (`v0.1.0`, not `0.1.0`).
- Pre-1.0 (`0.x.y`): breaking changes may bump MINOR; PATCH for fixes and small features.

## Before releasing

1. Ensure CI is green on `main`.
2. Update `CHANGELOG.md`:
   - Move items from `[Unreleased]` into a new section: `## [x.y.z] - YYYY-MM-DD`.
   - Leave an empty `[Unreleased]` section at the top.
3. Commit the changelog update (message e.g. `chore: release v0.1.0`).

## Tag the release

Tag the changelog commit **before** pushing or publishing the GitHub release:

```bash
VERSION=v0.1.0

git tag -a "$VERSION" -m "$VERSION"
```

Verify the tag points at the right commit: `git show "$VERSION" --no-patch`.

## Create the release

```bash
VERSION=v0.1.0
ASSET=cards-linux-amd64

# Push commit and tag
git push origin main
git push origin "$VERSION"

# Build release binary (linux/amd64) from the tagged commit
LDFLAGS="-s -w \
  -X github.com/joaovictornsv/cards-cli/internal/buildinfo.Version=${VERSION#v} \
  -X github.com/joaovictornsv/cards-cli/internal/buildinfo.Commit=$(git rev-parse HEAD)"
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o "$ASSET" ./cmd/cards

# Generate checksum file for the release binary
sha256sum "$ASSET" > SHA256SUMS

# Create GitHub release with binary and checksum attached (notes: see "Release notes format" below)
gh release create "$VERSION" \
  --title "$VERSION" \
  --notes "$(cat <<'EOF'
- <opening line: one-sentence theme for this release>
- <bullet: user-visible change>
- <bullet: user-visible change>
EOF
)" \
  "$ASSET" SHA256SUMS

rm -f "$ASSET" SHA256SUMS
```

Every release must include the built binary and a `SHA256SUMS` file. Do not publish a release without uploading both `cards-linux-amd64` and `SHA256SUMS`.

Users verify downloaded binaries with:

```bash
sha256sum -c SHA256SUMS
```

Add extra GOOS/GOARCH builds only when explicitly requested; default is linux/amd64 only.

## Release notes format

Read the version section in `CHANGELOG.md`, then write **concise bullet-point notes** for the GitHub release. Do not paste the raw changelog or use a generic placeholder like "See CHANGELOG.md for details."

### Rules

- **Bullet list only** — every line starts with `-`.
- **3–5 bullets** for most releases; split Added/Removed into separate bullets when both exists.
- **First bullet** — one short sentence framing the release (theme or scope).
- **Following bullets** — user-visible changes only; group related items in one bullet when it stays readable.
- **Concise, not exhaustive** — summarize smartly; skip internal/CI/docs-only items unless they matter to users.
- **No extra sections** — no headings, no "See CHANGELOG.md" footer.

### Editing an existing release

```bash
gh release edit v0.1.0 --notes "$(cat <<'EOF'
- ...
EOF
)"
```

## Checklist

- [ ] CHANGELOG updated for this version
- [ ] Annotated tag created locally (`git tag -a vX.Y.Z -m vX.Y.Z`)
- [ ] Commit and tag pushed to `origin`
- [ ] `cards-linux-amd64` built from the tagged commit
- [ ] `SHA256SUMS` generated from the release binary
- [ ] Release notes written as concise bullets (see format above)
- [ ] GitHub release created with matching notes, binary, and `SHA256SUMS` uploaded
