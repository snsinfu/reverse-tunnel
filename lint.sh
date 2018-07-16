#!/bin/sh

run() {
    printf "\033[36m%s\033[m\n" "$*"
    printf "\033[33m"
    "$@"
    printf "\033[m"
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
