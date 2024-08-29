#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo "Must provide .env file as argument"
    echo "Usage: $0 <path_to_env_file>"
    exit 1
fi

if [ -f $1 ]; then
    export $(grep -v '^#' $1 | xargs) 
else
    echo "$1 not found"
    exit 1
fi
