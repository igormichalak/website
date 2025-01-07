#!/bin/bash

go mod download
go build -trimpath -ldflags "-s -w" -o ./bin/httpserver ./cmd/httpserver/
