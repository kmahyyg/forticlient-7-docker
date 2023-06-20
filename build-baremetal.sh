#!/bin/bash
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

cd ./gosrc
git describe --long --dirty --tags --always | tr -d '\n' > ./bareMetalAnswerBot/version.txt
go build -o ./bin/go-fortivpn-daemon -trimpath -ldflags='-s -w' ./bareMetalAnswerBot/main.go
cd ../
