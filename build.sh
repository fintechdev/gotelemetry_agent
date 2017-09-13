#!/bin/bash

if [ "x$ARTIFACTS_DIR" = "x" ] ; then
  echo "ERROR: the environment variable ARTIFIACTS_DIR must be set"
  exit 1
fi

# should be the standard for every script
OLD_DIR=$( pwd )
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR

# now our stuff
os=$(go env GOOS)
arch=$(go env GOARCH)

# Get the git commit
GIT_COMMIT="$(git rev-parse HEAD)"
GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

os=linux
arch=amd64
echo
echo "Building gotelemetry_agent-$os-$arch now..."
GOOS="$os" GOARCH="$arch" CGO_ENABLED=1 go build -v \
  -o $ARTIFACTS_DIR/gotelemetry_agent-$os-$arch \
  -ldflags "-X github.com/telemetryapp/gotelemetry_agent/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}'" \
  github.com/telemetryapp/gotelemetry_agent

os=darwin
arch=amd64
echo
echo "Building gotelemetry_agent-$os-$arch now..."
GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build -v \
  -o $ARTIFACTS_DIR/gotelemetry_agent-$os-$arch \
  -ldflags "-X github.com/telemetryapp/gotelemetry_agent/version.GitCommit='${GIT_COMMIT}${GIT_DIRTY}'" \
  github.com/telemetryapp/gotelemetry_agent

cd $OLD_DIR
