#!/bin/bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

git describe --long --dirty --tags --always | tr -d '\n' > ./gosrc/bareMetalAnswerBot/version.txt
go build -o ./bin/go-fortivpn-daemon -trimpath -ldflags='-s -w' ./gosrc/bareMetalAnswerBot/main.go
