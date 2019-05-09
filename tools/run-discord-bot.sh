#!/bin/bash
cd "$(dirname "$0")" || exit

source config.env
go run ../cmd/discord-bot/main.go
