#!/usr/bin/env bash

# Script to run go test on multiple packages with a single
# coverage report for codecov.io.
# You probably just want to use "go test ./..." instead of
# this script
set -e
echo "" > coverage.txt

for d in $(find ./* -maxdepth 10 -type d); do
    if ls $d/*.go &> /dev/null; then
        go test -coverprofile=profile.out -covermode=atomic $d
        if [ -f profile.out ]; then
            cat profile.out >> coverage.txt
            rm profile.out
        fi
    fi
done
