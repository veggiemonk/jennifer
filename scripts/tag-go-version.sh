#!/usr/bin/env bash
set -euo pipefail

# Creates a semver git tag from the Go version in go.mod.
# Example: go.mod says "go 1.26.1" → tag "v1.26.1"

GO_VERSION=$(grep '^go ' go.mod | awk '{print $2}')
if [[ -z "$GO_VERSION" ]]; then
    echo "error: could not parse Go version from go.mod" >&2
    exit 1
fi

TAG="v${GO_VERSION}"

if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "tag $TAG already exists"
    exit 0
fi

echo "Creating tag: $TAG"
git tag -s "$TAG" -m "Update to Go ${TAG}"
echo "Done. Push with: git push origin $TAG"
