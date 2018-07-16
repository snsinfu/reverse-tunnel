#!/bin/sh

run() {
    printf "[36m%s[m\n" "$*"
    printf "[33m"
    "$@"
    printf "[m"
}

gocyclo_check() {
    gocyclo "$@" | awk '$1 > 15'
}

run gofmt -l .
run go vet ./...
run ineffassign .
run misspell .
run gocyclo_check .
run golint ./...
run megacheck ./...
