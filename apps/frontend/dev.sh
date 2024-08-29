#!/usr/bin/env bash

source ./load_env.sh .env
cd server
go run ./cmd/core/main.go
