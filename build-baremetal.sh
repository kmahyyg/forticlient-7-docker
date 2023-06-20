#!/bin/bash
set -xe
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

cd ./gosrc
git describe --long --dirty --tags --always | tr -d '\n' > ./bareMetalAnswerBot/version.txt
go build -o ../fortivpn-autobot/usr/local/fortivpn_autobot/go-fortivpn-daemon -trimpath -ldflags='-s -w' ./bareMetalAnswerBot/main.go
cd ../fortivpn-autobot
