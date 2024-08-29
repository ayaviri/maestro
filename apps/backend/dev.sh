#!/usr/bin/env bash

function core {
    source ./load_env.sh .env.core
    go run ./cmd/core
}

function fs {
    source ./load_env.sh .env.fs
    go run ./cmd/fs/main.go
}

function worker {
    source ./load_env.sh .env.worker
    go run ./worker
}

$@
