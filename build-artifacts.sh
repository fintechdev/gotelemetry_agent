#!/bin/bash

# should be the standard for every script
OLD_DIR=$( pwd )
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

DOCKER=$( which docker )
if [ $? -ne 0 ] ; then
  echo "ERROR: you must have docker installed"
  exit 1
fi

rm -vf $DIR/artifacts/gotelemetry_agent-linux-amd64 || echo no previous artifact gotelemetry_agent-linux-amd64 found
rm -vf $DIR/artifacts/gotelemetry_agent-darwin-amd64 || echo no previous artifact gotelemetry_agent-darwin-amd64 found

cd $DIR
docker build -t gotelemetry_agent-artifacts-builder .
ret=$?
if [ $ret -ne 0 ] ; then
  cd $OLD_DIR
  exit $ret
fi

docker run \
  --rm \
  -v $DIR/artifacts:/artifacts \
  -e ARTIFACTS_DIR=/artifacts \
  gotelemetry_agent-artifacts-builder
ret=$?
if [ $ret -ne 0 ] ; then
  cd $OLD_DIR
  exit $ret
fi
cd $OLD_DIR
