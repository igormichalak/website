#!/bin/bash

go build -trimpath -ldflags "-s -w" -o ./bin/httpserver ./cmd/httpserver/
