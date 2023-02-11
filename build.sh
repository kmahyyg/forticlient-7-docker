#!/bin/bash
git describe --long --dirty --tags --always > ./gosrc/cmd/version.txt
sudo podman build . -t ghcr.io/kmahyyg/forticlient:7 -f ./Forticlient.dockerfile
