#!/bin/bash
git describe --long --dirty --tags --always | tr -d '\n' > ./gosrc/cmd/version.txt
sudo podman build . -t ghcr.io/kmahyyg/fortivpn:7 -f ./Forticlient.dockerfile
