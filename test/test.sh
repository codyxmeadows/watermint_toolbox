#!/usr/bin/env bash

TEST_RESULTS=./test/results
TEST_DEBUG=./test/debug
TEST_OUT=$TEST_DEBUG/all.out
TEST_ERR=$TEST_DEBUG/err.out
TEST_PROFILE=coverage.txt

mkdir -p resources/keys
touch resources/keys/toolbox.build
mkdir -p $TEST_RESULTS
mkdir -p $TEST_DEBUG

if [ -e test/target_package ]; then
  TEST_PACKAGES=$(cat test/target_package)
else
  TEST_PACKAGES=./...
fi
TEST_PACKAGES_SUM=$(echo $TEST_PACKAGES | shasum -a 256 | awk '{print $1}')

echo TEST: Run tests: $TEST_PACKAGES
go test -v -short -timeout 30s -covermode=atomic -coverprofile=$TEST_PROFILE $TEST_PACKAGES >"$TEST_OUT" 2>"$TEST_ERR"
TEST_EXIT_CODE=$?

if [ $TEST_EXIT_CODE -ne 0 ]; then
  echo Test failed: $TEST_EXIT_CODE
  if [ "$CIRCLE_BUILD_NUM"x != ""x ]; then
    echo TEST: Uploading logs
    find $TEST_DEBUG
    go run tbx.go dev ci artifact up -budget-memory low -local-path $TEST_DEBUG -dropbox-path /watermint-toolbox-build/test-logs/$CIRCLE_BUILD_NUM/$TEST_PACKAGES_SUM -peer-name deploy
    go run tbx.go dev ci artifact up -budget-memory low -local-path $HOME/.toolbox/jobs -dropbox-path /watermint-toolbox-build/test-logs/$CIRCLE_BUILD_NUM/$TEST_PACKAGES_SUM -peer-name deploy
  fi
  exit 1
fi

echo TEST: Finished
