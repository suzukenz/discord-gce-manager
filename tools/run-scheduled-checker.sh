#!/bin/bash
cd "$(dirname "$0")" || exit

source config.env
go run ../cmd/scheduled-checker/main.go
