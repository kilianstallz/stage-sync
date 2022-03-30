#!/bin/bash

set -eu

if ! git diff-index --quiet HEAD -- ; then
    echo "uncommited changes on HEAD, aborting"
    exit 1
fi

npx standard-version

VERSION=$(git describe --tags "$(git rev-list --tags --max-count=1)")

cat > internal/version.go <<EOF
package version

const Version = "$VERSION"
EOF

git add .
git commit -m "chore(prepare): $VERSION"

cat > internal/version.go <<EOF
package version

const Version = "$VERSION-dev"
EOF

git add .
git commit -m "$VERSION postrelease bump"
git push --follow-tags origin main


echo "Now go write some release notes!"
