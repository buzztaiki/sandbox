#!/bin/sh

go run . "$(az account show | jq '.id' -r)" "$@"
